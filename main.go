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

	contributors, err := gitRepo.GetContributors()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%+v", contributors)
}
