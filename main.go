package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

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

var (
	org = flag.String("org", "automata-devops-io", "organization to target in github")
)

func repoMan(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	flag.Parse()
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	client := github.NewClient(tokenClient)

	payload, err := github.ValidatePayload(r, []byte("0d2aed9c2cf2ca01d2660d3723d89f07c94f9f8e"))
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
	case *github.RepositoryEvent:
		if e.Action != nil && *e.Action == "created" {
			fileContent := []byte("# Auto_generated_file\n ***This file was automatically created using the githubAPi and webhooks.***\n If you'd like more information on how this file was created, the api doc [here](https://docs.github.com/en/rest/repos/contents#create-a-file) ")
			opts := &github.RepositoryContentFileOptions{
				Message:   github.String("Initial commit"),
				Content:   fileContent,
				Branch:    github.String("main"),
				Committer: &github.CommitAuthor{Name: github.String("Jeff Brimager"), Email: github.String("jbrimager@gmail.com")},
			}
			preq := &github.ProtectionRequest{
				EnforceAdmins: true,
				RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
					RequiredApprovingReviewCount: 2,
					DismissStaleReviews:          true,
				},
			}
			client.Repositories.UpdateBranchProtection(ctx, *org, *e.Repo.Name, "main", preq)
			time.Sleep(2 * time.Second)
			client.Repositories.AddAdminEnforcement(ctx, *org, *e.Repo.Name, "main")
			time.Sleep(2 * time.Second)
			client.Repositories.CreateFile(ctx, *org, *e.Repo.Name, "check_me_out.md", opts)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
