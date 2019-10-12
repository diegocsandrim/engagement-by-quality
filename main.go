package main

import (
	"fmt"
	"log"

	"./git"
)

func main() {
	namespace := "kelseyhightower"
	project := "envconfig"

	err := git.EnsureGitBaseDirExists()
	if err != nil {
		log.Panic(err)
	}

	err = git.Clone(namespace, project)
	if err != nil {
		log.Panic(err)
	}

	constributors, err := git.GetContributors(namespace, project)
	fmt.Printf("%v", constributors)
}
