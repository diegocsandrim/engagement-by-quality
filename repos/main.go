package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type RepositoriesResult struct {
	Items []struct {
		FullName string `json:"full_name"`
	} `json:"items"`
}

// Finds the top 1000 most forked repositories.
// Filters used:
// - Only repositories with Golang
// - Only repositories created from 2012-01-01 to 2015-01-01
// - Only repositories that had a push after 2018-01-01
// Uses env GH_TOKEN for authentication
//

func main() {
	token, _ := os.LookupEnv("GH_TOKEN")
	if token == "" {
		log.Fatal("GH_TOKEN environment variable is not set. Tip: create the token at https://github.com/settings/tokens and set the env var GH_TOKEN to its value")
	}

	repositoryNames, err := fetchRepositoryNames(token)
	if err != nil {
		log.Fatal(err)
	}

	output := strings.Join(repositoryNames, "\n")
	fmt.Println(output)
}

func fetchRepositoryNames(token string) ([]string, error) {
	urlFormat := "https://api.github.com/search/repositories?sort=%s&per_page=%d&page=%d&q=%s"
	sort := "fork"
	perPage := 100
	lastPage := 10
	query := "language:Go+created:2012-01-01..2015-01-01+pushed:2018-01-01..*+is:public"

	var repositoriesResult RepositoriesResult

	repositoryNames := make([]string, 0, 1000)

	for page := 1; page <= lastPage; page++ {
		log.Printf("Requesting page %d", page)
		url := fmt.Sprintf(urlFormat, sort, perPage, page, query)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create the request: %w", err)
		}

		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

		http.DefaultClient.Timeout = 1 * time.Minute

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page result: %w", err)
		}

		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&repositoriesResult)
		if err != nil {
			return nil, fmt.Errorf("failed decoding response: %w", err)
		}

		for _, repositoryResponse := range repositoriesResult.Items {
			repositoryNames = append(repositoryNames, repositoryResponse.FullName)
		}

	}

	return repositoryNames, nil
}
