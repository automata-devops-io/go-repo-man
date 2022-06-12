package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	port := os.Getenv("PORT")
	log.Println("server started")
	http.HandleFunc("/test", repoMan)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// var (
// 	// org name will be a
// 	org = flag.String("org", "automata-devops-io", "organization to target in github")
// )

func repoMan(w http.ResponseWriter, r *http.Request) {
	ghtoken := os.Getenv("GHTOKEN")
	whsecret := os.Getenv("WHSECRET")
	ctx := context.Background()
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghtoken},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	client := github.NewClient(tokenClient)

	payload, err := github.ValidatePayload(r, []byte(whsecret))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		// this is a commit push, do something with it
	case *github.PullRequestEvent:
		// this is a pull request, do something with it
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
		// this is for setting up Repo security
	case *github.RepositoryEvent:
		if e.Action != nil && *e.Action == "created" {
			issue := &github.IssueRequest{
				Title:    github.String("New repo Created"),
				Body:     github.String("@sam1el this repo was created"),
				Assignee: github.String("sam1el"),
			}
			preq := &github.ProtectionRequest{
				EnforceAdmins: true,
				RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
					RequiredApprovingReviewCount: 2,
					DismissStaleReviews:          true,
					RequireCodeOwnerReviews:      true,
				},
			}
			client.Repositories.UpdateBranchProtection(ctx, *e.Org.Name, *e.Repo.Name, "main", preq)
			client.Repositories.AddAdminEnforcement(ctx, *e.Org.Name, *e.Repo.Name, "main")
			client.Issues.Create(ctx, *e.Org.Name, *e.Repo.Name, issue)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
