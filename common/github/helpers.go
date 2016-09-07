package github

import (
	"fmt"

	"github.com/skeswa/gophr/common/models"
)

// ParseStarCount TODO Won't need this after implementing FFJSON
func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}

// BuildGitHubBranch creates a new ref based on a hash of the old ref
func BuildGitHubBranch(ref string) string {
	repoHash := ref[:len(ref)-1]
	return repoHash
}

// BuildRemoteURL creates a remote url for a packageModel based on it's ref
func BuildRemoteURL(packageModel *models.PackageModel, ref string) string {
	repoName := BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	remoteURL := fmt.Sprintf(gitHubRemoteOrigin, repoName)
	return remoteURL
}

// BuildNewGitHubRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similiar usernames and packages
func BuildNewGitHubRepoName(author string, repo string) string {
	return fmt.Sprintf("%d%s%d%s", len(author), author, len(repo), repo)
}