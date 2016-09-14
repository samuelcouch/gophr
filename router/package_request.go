package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
	"github.com/skeswa/gophr/common/subv"
)

const (
	formKeyGoGet                = "go-get"
	formValueGoGet              = "1"
	gophrDomainDev              = "gophr.dev"
	gophrDomainProd             = "gophr.prod"
	contentTypeHTML             = "text/html"
	subPathRegexStr             = `((?:\/[a-zA-Z0-9][-.a-zA-Z0-9]*)*)`
	userRepoRegexStr            = `^\/([a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\/([a-zA-Z0-9\.\-_]+)`
	masterGitRefLabel           = "master"
	someFakeGitTagRef           = "refs/tags/thisisnotathinginanyrepo"
	gitRefsInfoSubPath          = "/info/refs"
	httpLocationHeader          = "Location"
	refSelectorRegexStr         = "([a-fA-F0-9]{40})"
	gitUploadPackSubPath        = "/git-upload-pack"
	httpContentTypeHeader       = "Content-Type"
	packagePageURLTemplate      = "https://%s/#/packages/%s/%s"
	contentTypeGitUploadPack    = "application/x-git-upload-pack-advertisement"
	githubUploadPackURLTemplate = "https://github.com/%s/%s/git-upload-pack"
	// packageRequestRegexTemplate:
	// "/author/repo@semver" or ""/author/repo@semver/subpath"
	packageRequestRegexTemplate = `%s(?:@%s)%s$`
	// versionSelectorRegexTemplate:
	//
	versionSelectorRegexTemplate = `([\%c\%c]?)([0-9]+)(?:\.([0-9]+|%c))?(?:\.([0-9]+|%c))?(?:\-([a-zA-Z0-9\-_]+[a-zA-Z0-9])(?:\.([0-9]+|%c))?)?([\%c\%c]?)`
	// barePackageRequestRegexTemplate:
	// "/author/repo" or ""/author/repo/subpath"
	barePackageRequestRegexTemplate = `%s%s$`
	masterRefName                   = "refs/heads/master"
)

// PackageRequest is stuct that standardizes the output of all the scenarios
// through which a package may be requested. PackageRequest is essentially a
// helper struct to move data between the sub-functions of
// RespondToPackageRequest and RespondToPackageRequest itself.
type packageRequest struct {
	req             *http.Request
	parts           *packageRequestParts
	refsData        []byte
	matchedSHA      string
	matchedSHALabel string
}

// newPackageRequestArgs is the arguments struct for newPackageRequest.
type newPackageRequestArgs struct {
	req        *http.Request
	downloader refsDownloader
}

// newPackageRequest parses and simplifies the information in a package version
// request in order to make serializing a response easier.
func newPackageRequest(args newPackageRequestArgs) (*packageRequest, error) {
	// Read the parts of the package request.
	parts, err := readPackageRequestParts(args.req)
	if err != nil {
		return nil, err
	}

	var (
		refs            common.Refs
		matchedSHA      string
		matchedSHALabel string
		packageRefsData []byte
	)

	// Only go out to fetch refs if they're going to get used.
	if isGoGetRequest(args.req) || isInfoRefsRequest(parts) {
		// Get and process all of the refs for this package.
		if refs, err = args.downloader.downloadRefs(
			parts.author,
			parts.repo); err != nil {
			return nil, err
		}

		// Set the default matched sha.
		matchedSHA = refs.MasterRefHash

		// If there are no candidates, return in failure.
		if refs.Candidates == nil || len(refs.Candidates) < 1 {
			return nil, NewNoSuchPackageVersionError(
				parts.author,
				parts.repo,
				parts.semverSelector.String())
		}

		// Figure out what the best candidate is.
		if parts.hasSemverSelector() {
			// Find the best candidate.
			bestCandidate := refs.Candidates.Best(parts.semverSelector)
			// Re-serialize the refs data with said candidate.
			matchedSHA = bestCandidate.GitRefHash
			matchedSHALabel = bestCandidate.GitRefLabel
			packageRefsData = refs.Reserialize(
				bestCandidate.GitRefName,
				bestCandidate.GitRefHash)
		} else if parts.hasSHASelector() {
			// Re-serialize the refs data with the sha.
			matchedSHA = parts.shaSelector
			packageRefsData = refs.Reserialize(
				someFakeGitTagRef,
				parts.shaSelector)
			// TODO(skeswa): investigate validating the ref to see if it actually
			// exists.
		} else {
			// Since there was no selector, we are fine with the fact that we didn't
			// find a match. Now, return the original refs that we downloaded from
			// github that point to master by default.
			packageRefsData = refs.Data
		}
	}

	return &packageRequest{
		req:             args.req,
		parts:           parts,
		refsData:        packageRefsData,
		matchedSHA:      matchedSHA,
		matchedSHALabel: matchedSHALabel,
	}, nil
}

// respondToPackageRequestArgs is the arguments struct for
// packageRequest#respond.
type respondToPackageRequestArgs struct {
	res     http.ResponseWriter
	conf    *config.Config
	creds   *config.Credentials
	session *gocql.Session
}

// respond crafts an appropriate response for a package request, serializes the
// aforesaid response and sends it back to the original client.
func (pr *packageRequest) respond(args respondToPackageRequestArgs) error {
	// Take care of the cases that deoend inf variations in the subpath.
	switch pr.parts.subpath {
	case gitUploadPackSubPath:
		// Send a 301 stipulating the repository can be found on github.
		args.res.Header().Set(
			httpLocationHeader,
			fmt.Sprintf(
				githubUploadPackURLTemplate,
				pr.parts.author,
				pr.parts.repo))
		args.res.WriteHeader(http.StatusMovedPermanently)
		return nil
	case gitRefsInfoSubPath:
		// Return the adjusted refs data when refs info is requested.
		args.res.Header().Set(httpContentTypeHeader, contentTypeGitUploadPack)
		args.res.Write(pr.refsData)
		return nil
	}

	// This means that go-get is requesting package/repository metadata.
	if isGoGetRequest(pr.req) {
		// Without blocking, count go-get surveying this package for installation as
		// a download in the database.
		go recordDownload(
			args.session,
			pr.parts.author,
			pr.parts.repo,
			// TODO(skeswa): record version label as well.
			// pr.matchedSHALabel,
			pr.matchedSHA)

		// Only run the sub-versioning if its completely necessary.
		if !models.IsPackageArchived(
			args.session,
			pr.parts.author,
			pr.parts.repo,
			pr.matchedSHA) {
			// Indicate in the logs that the package was archived.
			log.Printf(
				"Package %s/%s@%s has not yet been archived.\n",
				pr.parts.author,
				pr.parts.repo,
				pr.matchedSHA)

			// Perform sub-versioning.
			if err := subv.SubVersionPackageModel(
				args.conf,
				args.session,
				args.creds,
				&models.PackageModel{Author: &pr.parts.author, Repo: &pr.parts.repo},
				pr.matchedSHA); err != nil {
				// Report the sub-versioning failure to the logs.
				log.Printf(
					"sub-versioning failed for package %s/%s@%s: %v\n",
					pr.parts.author,
					pr.parts.repo,
					pr.matchedSHA,
					err)

				return err
			}
		}

		// Change the domain depending on whether this is dev or not.
		var domain string
		if args.conf.IsDev {
			domain = gophrDomainDev
		} else {
			domain = gophrDomainProd
		}

		// Compile the go-get metadata accordingly.
		var (
			repo     = github.BuildNewGitHubRepoName(pr.parts.author, pr.parts.repo)
			author   = github.GitHubGophrPackageOrgName
			metaData = []byte(generateGoGetMetadata(generateGoGetMetadataArgs{
				gophrURL:        generateGophrURL(domain, author, repo, pr.matchedSHA),
				treeURLTemplate: generateGithubTreeURLTemplate(author, repo, pr.matchedSHA),
				blobURLTemplate: generateGithubBlobURLTemplate(author, repo, pr.matchedSHA),
			}))
		)

		// Return the go-get metadata.
		args.res.Header().Set(httpContentTypeHeader, contentTypeHTML)
		args.res.Write(metaData)
		return nil
	}

	// If none of the other cases matched, then redirect to the package page.
	// TODO(skeswa): make this redirect specific to the version of the package.
	args.res.Header().Set(
		httpLocationHeader,
		fmt.Sprintf(
			packagePageURLTemplate,
			pr.req.URL.Host,
			pr.parts.author,
			pr.parts.repo))
	args.res.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// isGoGetRequest returns true if the request was made by go get.
func isGoGetRequest(req *http.Request) bool {
	return req.FormValue(formKeyGoGet) == formValueGoGet
}

// isInfoRefsRequest returns true if the request parts reflect that the request
// is a git refs info request.
func isInfoRefsRequest(parts *packageRequestParts) bool {
	return parts.subpath == gitRefsInfoSubPath
}

// recordDownload is a helper function that records the download of a specific
// package.
func recordDownload(
	session *gocql.Session,
	author string,
	repo string,
	selector string) {
	if err := models.RecordDailyDownload(
		session,
		author,
		repo,
		selector); err != nil {
		// Instead of bubbling this error, just commit it to the logs. That way this
		// failure is allowed to remain low impact.
		log.Printf(
			"Failed to record download for package %s/%s@%s: %v\n",
			author,
			repo,
			selector,
			err,
		)
	}
}
