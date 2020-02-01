package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
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
		"ac05a422f5d81b015e01f1a1e01a1344c542ea2a",
		"http://localhost:9000",
		gitRepo.ProjectDir(),
	)

	day := time.Hour * 24
	analysisDate := time.Now().UTC().Add(-day * time.Duration(len(contributorAttractorCommits)))
	reader := bufio.NewReader(os.Stdin)
	silent := false
	for _, contributorAttractorCommit := range contributorAttractorCommits {
		analysisDate = analysisDate.Add(day)
		shortCommitHash := contributorAttractorCommit.Commit.Hash[0:8]
		log.Printf("Analysing commit %s\n", shortCommitHash)
		err = gitRepo.Checkout(contributorAttractorCommit.Commit.Hash)
		if err != nil {
			log.Panic(err)
		}

		for {
			if silent {
				break
			}

			fmt.Print("Continue Y (yes), n (no), s (silent)?: ")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(strings.ToLower(text), "\n")
			fmt.Println(text)

			if text == "y" || text == "" {
				break
			}

			if text == "n" {
				os.Exit(0)
			}

			if text == "s" {
				silent = true
				break
			}

		}

		err = qualityAnalyzer.Run(shortCommitHash, analysisDate)
		if err != nil {
			log.Panic(err)
		}
	}
	// TODO (minor): Remover arquivos temporários deixados pelo sonnar-scanner (problema de lock e de ownership)
	// Ideia: docker exec usando o sonnar-scanner depois que ele terminou (nem sei se é possível)
}
