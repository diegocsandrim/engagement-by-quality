package main

import (
	"log"

	"./git"
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

	for _, contributor := range gitRepo.Contributors() {
		log.Printf("%s %s %s", contributor.Id, contributor.FirstCommit().Hash, contributor.FirstCommit().ParentHash)
	}
}
