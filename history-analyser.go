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

func analyseHistory(namespace string, project string, sonarKey string, sonarUrl string, iterative bool) error {
	gitRepo := git.NewGitRepo(namespace, project)

	err := gitRepo.Clone()
	if err != nil {
		return fmt.Errorf("could not clone repo: %w", err)
	}

	err = gitRepo.LoadCommits()
	if err != nil {
		return fmt.Errorf("could not load commits: %w", err)
	}

	contributorAttractorCommits := gitRepo.ContributorAttractorCommits()
	sort.Slice(contributorAttractorCommits, func(i, j int) bool {
		commitI := contributorAttractorCommits[i].Commit
		commitJ := contributorAttractorCommits[j].Commit
		return commitI.Date.Before(commitJ.Date)
	})

	qualityAnalyzer := qualityanalyzers.NewSonnar(
		qualityanalyzers.FormatProjectKey(namespace, project),
		sonarKey,
		sonarUrl,
		gitRepo.ProjectDir(),
	)

	day := time.Hour * 24
	analysisDate := time.Now().UTC().Add(-day * time.Duration(len(contributorAttractorCommits)))
	reader := bufio.NewReader(os.Stdin)

	for i, contributorAttractorCommit := range contributorAttractorCommits {
		analysisDate = analysisDate.Add(day)
		shortCommitHash := contributorAttractorCommit.Commit.Hash[0:8]
		log.Printf("Analysing commit %s (%d/%d)\n", shortCommitHash, i+1, len(contributorAttractorCommits))
		err = gitRepo.Checkout(contributorAttractorCommit.Commit.Hash)
		if err != nil {
			return fmt.Errorf("could not checkout to commit: %w", err)
		}

		for {
			if !iterative {
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
				iterative = false
				break
			}

		}
		if dryRun() {
			break
		}
		err = qualityAnalyzer.Run(shortCommitHash, analysisDate)
		if err != nil {
			return fmt.Errorf("could not run analyser: %w", err)
		}
	}
	return nil
	// TODO (minor): Remover arquivos temporários deixados pelo sonnar-scanner (problema de lock e de ownership)
	// Ideia: docker exec usando o sonnar-scanner depois que ele terminou (nem sei se é possível)
}

func dryRun() bool {
	return false
}
