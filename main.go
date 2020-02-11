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
		// {namespace: "kelseyhightower", project: "envconfig"},
		{namespace: "fatedier", project: "frp"},
		{namespace: "keybase", project: "client"},
		{namespace: "syncthing", project: "syncthing"},
		{namespace: "gin-gonic", project: "gin"},
		{namespace: "cockroachdb", project: "cockroach"},
		{namespace: "openshift", project: "origin"},
		{namespace: "pingcap", project: "tidb"},
		{namespace: "gogs", project: "gogs"},
		{namespace: "prometheus", project: "prometheus"},
		{namespace: "istio", project: "istio"},
		{namespace: "etcd-io", project: "etcd"},
		{namespace: "gohugoio", project: "hugo"},
		{namespace: "kubernetes/test", project: "infra"},
		{namespace: "Alluxio", project: "alluxio"},
		{namespace: "hashicorp", project: "terraform"},
		{namespace: "golang", project: "go"},
		{namespace: "moby", project: "moby"},
		{namespace: "kubernetes", project: "kubernetes"},
		{namespace: "helm", project: "charts"},
	}

	sonarKey := "390013c5cbfe8ece1c357436cf54402336ad1d46"
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
