package main

import (
	"log"
)

func main() {
	// Some test list from https://github.com/search?l=&o=desc&q=language%3AGo+stars%3A%22%3E+15000%22&s=stars&type=Repositories
	githubProjects := []struct {
		namespace string
		project   string
	}{
		{namespace: "Alluxio", project: "alluxio"},

		// {namespace: "kelseyhightower", project: "envconfig"},
		// {namespace: "keybase", project: "client"},
		// {namespace: "helm", project: "charts"},
		// {namespace: "cockroachdb", project: "cockroach"},
		// {namespace: "kubernetes", project: "test-infra"},
		// {namespace: "openshift", project: "origin"},
		// {namespace: "pingcap", project: "tidb"},
		// {namespace: "kubernetes", project: "kubernetes"},
		// {namespace: "hashicorp", project: "terraform"},
		// {namespace: "istio", project: "istio"},
		// {namespace: "moby", project: "moby"},
		// {namespace: "golang", project: "go"},
		// {namespace: "kubernetes", project: "kubernetes"},
		// {namespace: "moby", project: "moby"},
		// {namespace: "gohugoio", project: "hugo"},
		// {namespace: "gin-gonic", project: "gin"},
		// {namespace: "gogs", project: "gogs"},
		// {namespace: "fatedier", project: "frp"},
		// {namespace: "syncthing", project: "syncthing"},
		// {namespace: "etcd-io", project: "etcd"},
		// {namespace: "prometheus", project: "prometheus"},
	}

	sonarKey := "0971341b9b41282b6216bb222bc326245098a822"
	sonarUrl := "http://localhost:9000"
	iterative := false

	for _, githubProject := range githubProjects {
		log.Printf("start project %s/%s", githubProject.namespace, githubProject.project)
		err := analyseHistory(githubProject.namespace, githubProject.project, sonarKey, sonarUrl, iterative)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("finished project %s/%s", githubProject.namespace, githubProject.project)
	}
}
