package main

import (
	"log"
	"sort"

	"./git"
	"./qualityanalyzers"
)

func main() {
	namespace := "kelseyhightower"
	project := "envconfig"

	gitRepo := git.NewGitRepo(namespace, project)

	err := gitRepo.ForceClone()
	if err != nil {
		log.Panic(err)
	}

	err = gitRepo.LoadCommits()
	if err != nil {
		log.Panic(err)
	}

	contributorAttractorCommits := gitRepo.ContributorAttractorCommits()
	sort.Slice(contributorAttractorCommits, func(i, j int) bool {
		commitI := contributorAttractorCommits[i].Commit
		commitJ := contributorAttractorCommits[j].Commit
		return commitI.Date.Before(commitJ.Date)
	})

	qualityAnalyzer := qualityanalyzers.NewSonnar(
		qualityanalyzers.FormatProjectKey(namespace, project),
		"e08c00b8a357612588fecabe92d1eb9971c7b74b",
		"http://localhost:9000",
		gitRepo.ProjectDir(),
	)

	for _, contributorAttractorCommit := range contributorAttractorCommits {
		log.Printf("Analysing commit %s\n", contributorAttractorCommit.Commit.Hash)
		err = gitRepo.Checkout(contributorAttractorCommit.Commit.Hash)
		if err != nil {
			log.Panic(err)
		}

		err = qualityAnalyzer.Run(contributorAttractorCommit.Commit.Hash[:8], contributorAttractorCommit.Commit.Date)
		if err != nil {
			log.Panic(err)
		}
	}
}
