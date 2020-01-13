package main

import (
	"log"
	"sort"
	"time"

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

	day := time.Hour * 24
	analysisDate := time.Now().UTC().Add(-day * time.Duration(len(contributorAttractorCommits)))
	for _, contributorAttractorCommit := range contributorAttractorCommits {
		analysisDate = analysisDate.Add(day)
		shortCommitHash := contributorAttractorCommit.Commit.Hash[0:8]
		log.Printf("Analysing commit %s\n", shortCommitHash)
		err = gitRepo.Checkout(contributorAttractorCommit.Commit.Hash)
		if err != nil {
			log.Panic(err)
		}

		err = qualityAnalyzer.Run(shortCommitHash, analysisDate)
		if err != nil {
			log.Panic(err)
		}
	}
	// TODO (minor): Remover arquivos temporários deixados pelo sonnar-scanner (problema de lock e de ownership)
	// Ideia: docker exec usando o sonnar-scanner depois que ele terminou (nem sei se é possível)
}
