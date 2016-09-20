package main

import (
	"fmt"
	"log"
	"time"

	"github.com/skeswa/gophr/common/verdeps"
)

const (
	// archiveExistenceCheckDelayMS is the time gap between existence check
	// attempts.
	archiveExistenceCheckDelayMS = 500
	// archiveExistenceCheckAttemptsLimit sets the cap on how many times an archive
	// existence check is attempted before the an error is recorded.
	archiveExistenceCheckAttemptsLimit = 3
)

// versionAndArchivePackage takes a package, locks all of its versions
// in a chronologically accurate way, and archives it in depot to be queried
// later.
func versionAndArchivePackage(args packageVersionerArgs) error {
	log.Printf("Preparing to sub-version %s/%s@%s \n", args.author, args.repo, args.sha)

	// Download the package in the construction zone.
	downloadPaths, err := args.downloadPackage(packageDownloaderArgs{
		sha:                  args.sha,
		repo:                 args.repo,
		author:               args.author,
		constructionZonePath: args.constructionZonePath,
	})
	if err != nil {
		return err
	}

	// Perform clean-up after function exits.
	defer args.attemptWorkDirDeletion(downloadPaths.workDirPath)

	// Version lock all of the Github dependencies in the packageModel.
	if err = args.versionDeps(verdeps.VersionDepsArgs{
		SHA:           args.sha,
		Repo:          args.repo,
		Path:          downloadPaths.archiveDirPath,
		Author:        args.author,
		GithubService: args.ghSvc,
	}); err != nil {
		return fmt.Errorf("Could not version deps properly: %v.", err)
	}

	// Create a new repository in the depot before pushing to it.
	if repoIsNew, repoCreationErr := args.createDepotRepo(
		args.author,
		args.repo,
		args.sha,
	); repoCreationErr != nil {
		return repoCreationErr
	} else if !repoIsNew {
		// If the repo is not new, that means this package is already being
		// versioned, or has already been versioned. So, we must wait for this
		// package to be archived.
		for attempts := 0; attempts < archiveExistenceCheckAttemptsLimit; attempts = attempts + 1 {
			// Enforce a time delay between attempts so as to allow for archival to
			// occur.
			if attempts > 0 {
				time.Sleep(archiveExistenceCheckDelayMS * time.Millisecond)
			}

			if archived, archiveCheckErr := args.isPackageArchived(packageArchivalArgs{
				db:     args.db,
				sha:    args.sha,
				repo:   args.repo,
				author: args.author,
			}); archiveCheckErr != nil {
				return fmt.Errorf(
					"Could not check if package has been versioned in another context: %v.",
					archiveCheckErr)
			} else if archived {
				// Since the package was archived elsewhere, exit here.
				return nil
			}
		}

		// The other package versioner context failed to deliver on time - complain
		// about it.
		return fmt.Errorf(
			"Was waiting for package archival of \"%s/%s@%s\" in another context, but it did not happen fast enough.",
			args.author,
			args.repo,
			args.sha)
	}

	// Push versioned package to depot, then delete the package directory from
	// the construction zone.
	if err = args.pushToDepot(packagePusherArgs{
		author:       args.author,
		repo:         args.repo,
		sha:          args.sha,
		creds:        args.creds,
		packagePaths: downloadPaths,
	}); err != nil {
		// Yikes, we couldn't push. So as to not prevent this package from ever
		// being versioned correctly, undo all the work we just did.
		if deletionError := args.destroyDepotRepo(
			args.author,
			args.repo,
			args.sha,
		); deletionError != nil {
			// This is wayy the worst case scenario here.
			return fmt.Errorf(
				"Could not delete package package in depot: %v. Had to delete because of a push error: %v.",
				deletionError,
				err)
		}

		return fmt.Errorf("Could not push versioned package to depot: %v.", err)
	}

	// Record that this package has been archived.
	go args.recordPackageArchival(packageArchivalArgs{
		db:     args.db,
		sha:    args.sha,
		repo:   args.repo,
		author: args.author,
	})

	return nil
}
