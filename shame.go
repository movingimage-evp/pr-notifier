package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var httpClient http.Client
var githubRepo = os.Getenv("GITHUB_REPOSITORY")
var githubToken = os.Getenv("GITHUB_TOKEN")

type GitHubPullRequest struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Number    int64     `json:"number"`
}

func githubRequest(request *http.Request) (*http.Response, error) {
	request.Header.Set("Accept", "application/vnd.github.shadow-cat-preview+json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func main() {
	now := time.Now()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient = http.Client{Transport: customTransport, Timeout: time.Minute}
	log.Printf("listing PRs for repo %s\n", githubRepo)
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls", githubRepo), nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := githubRequest(request)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var prList []GitHubPullRequest
	if err = json.Unmarshal(body, &prList); err != nil {
		log.Fatal(err)
	}
	for _, pr := range prList {
		if pr.CreatedAt.Before(now.AddDate(0, 0, -2)) {

			// todo: add slack message handling here
			//age := time.Since(pr.CreatedAt)
		}
	}
}
