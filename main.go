package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var httpClient http.Client
var githubRepo = os.Getenv("GITHUB_REPOSITORY")
var githubToken = os.Getenv("GITHUB_TOKEN")
var daysBefore = os.Getenv("DAYS_BEFORE")
var slackApiKey = os.Getenv("SLACK_API_KEY")
var slackChannelId = os.Getenv("SLACK_CHANNEL_ID")

type GitHubPullRequest struct {
	CreatedAt time.Time            `json:"created_at"`
	Title     string               `json:"title"`
	Number    int64                `json:"number"`
	Reviewers []RequestedReviewers `json:"requested_reviewers"`
}

func (pr *GitHubPullRequest) hasPendingReviewers() bool {
	return len(pr.Reviewers) != 0
}

type RequestedReviewers struct {
	Login string
}

func ghRequest(request *http.Request) (*http.Response, error) {
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
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.github.com/repos/%s/pulls", githubRepo), nil)

	if err != nil {
		log.Fatal(err)
	}

	response, err := ghRequest(request)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("GitHub api request failed with status code %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var prList []GitHubPullRequest

	if err = json.Unmarshal(body, &prList); err != nil {
		log.Fatal(err)
	}

	daysBeforeConverted, err := strconv.Atoi(daysBefore)

	if err != nil {
		log.Fatal(err)
	}

	slackClient := slack.New(slackApiKey)
	ctx := context.Background()
	block := []slack.Block{slack.NewContextBlock("context", slack.NewTextBlockObject(slack.PlainTextType, "hello", false, false))}

	messageOption := []slack.MsgOption{
		slack.MsgOptionBlocks(block...),
		slack.MsgOptionText("text", false),
	}

	_, _, err = slackClient.PostMessageContext(ctx, slackChannelId, messageOption...)

	if err != nil {
		log.Fatal(err)
	}

	for _, pr := range prList {
		if pr.CreatedAt.Before(now.AddDate(0, 0, daysBeforeConverted)) {
			if pr.hasPendingReviewers() {
				log.Printf("PR %s has open reviewers %s\n", pr.Title, pr.Reviewers)
			}
		}
	}
}
