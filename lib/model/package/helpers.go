package pkg

import (
	"bytes"
	"sync"

	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model/package/awesome"
)

// checkIfAwesomeAsynchronously is a wrapper around awesome.IncludesPackage that makes it
// easier to place in a go-routine by making the return values pointers instead.
func checkIfAwesomeAsynchronously(
	q query.Queryable,
	author string,
	repo string,
	outputBool *bool,
	outputError *error,
	wg *sync.WaitGroup,
) {
	awesome, err := awesome.IncludesPackage(q, author, repo)
	if err != nil {
		*outputError = err
		wg.Done()
		return
	}

	*outputBool = awesome
	wg.Done()
}

// getGithubRepoDataAsynchronously is a wrapper around
// github.RequestService.FetchGitHubDataForPackageModel that makes it easier to
// place in a go-routine by making the return values pointers instead.
func getGithubRepoDataAsynchronously(
	ghSvc github.RequestService,
	author string,
	repo string,
	outputRepoData *dtos.GithubRepoDTO,
	outputError *error,
	wg *sync.WaitGroup,
) {
	repoData, err := ghSvc.FetchGitHubDataForPackageModel(author, repo)
	if err != nil {
		*outputError = err
		wg.Done()
		return
	}

	*outputRepoData = repoData
	wg.Done()
}

// composeSearchBlob uses other package metadata to create the string that is
// indexed for search purposes.
func composeSearchBlob(author, repo, description string) string {
	var buffer bytes.Buffer

	buffer.WriteString(author)
	buffer.WriteByte(' ')
	buffer.WriteString(repo)
	buffer.WriteByte(' ')
	buffer.WriteString(description)

	return buffer.String()
}