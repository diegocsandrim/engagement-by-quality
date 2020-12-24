package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	config := Config{}

	app := &cli.App{
		Usage: "analyse the project history using sonnar-scanner",
		Commands: []*cli.Command{
			{
				Name:    "analyse",
				Aliases: []string{"a"},
				Usage:   "analyse the repository history",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "sonarkey",
						Usage:       "Sonarqube token",
						EnvVars:     []string{"EQ_SONAR_TOKEN"},
						Required:    true,
						Destination: &(config.SonarKey),
					},
					&cli.StringFlag{
						Name:        "sonarurl",
						Usage:       "Sonarqube URL",
						EnvVars:     []string{"EQ_SONAR_URL"},
						Required:    true,
						Destination: &(config.SonarURL),
					},
					&cli.StringSliceFlag{
						Name:    "plugins",
						Usage:   "List of plugins to run with sonnar-scanner",
						EnvVars: []string{"EQ_SONAR_PLUGINS"},
					},
				},
				Action: func(c *cli.Context) error {
					config.SonarPlugins = c.StringSlice("plugins")

					if c.Args().Len() != 1 {
						return cli.Exit("argument must be exactly 1 repository", 1)
					}

					repositoryParts := strings.Split(c.Args().First(), "/")
					if len(repositoryParts) != 2 {
						return cli.Exit("argument must be in format namespace/project", 1)
					}

					namespace := repositoryParts[0]
					project := repositoryParts[1]

					log.Printf("starting namespace %s, project %s", namespace, project)

					err := analyseByMonth(namespace, project, config)
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					log.Printf("finished namespace %s, project %s", namespace, project)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func someProjects() {
	// Some test list from https://github.com/search?l=&o=desc&q=language%3AGo+stars%3A%22%3E+15000%22&s=stars&type=Repositories
	githubProjects := []struct {
		namespace string
		project   string
	}{
		{namespace: "kelseyhightower", project: "envconfig"},
		{namespace: "fatedier", project: "frp"},
		{namespace: "keybase", project: "client"},
		{namespace: "syncthing", project: "syncthing"},
		{namespace: "gin-gonic", project: "gin"},
		{namespace: "cockroachdb", project: "cockroach"},
		{namespace: "helm", project: "charts"},
		{namespace: "Alluxio", project: "alluxio"},
		{namespace: "gogs", project: "gogs"},
		{namespace: "openshift", project: "origin"},
		{namespace: "pingcap", project: "tidb"},
		{namespace: "prometheus", project: "prometheus"},
		{namespace: "istio", project: "istio"},
		{namespace: "etcd-io", project: "etcd"},
		{namespace: "gohugoio", project: "hugo"},
		{namespace: "kubernetes", project: "test-infra"},
		{namespace: "hashicorp", project: "terraform"},
		{namespace: "golang", project: "go"},
		{namespace: "moby", project: "moby"},
		{namespace: "kubernetes", project: "kubernetes"},
	}

	for _, githubProject := range githubProjects {
		log.Printf("start project %s/%s", githubProject.namespace, githubProject.project)
	}
}
