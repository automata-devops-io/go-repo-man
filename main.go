package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	port := os.Getenv("PORT")
	log.Info("server started")
	http.HandleFunc("/test", repoMan)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var (
	org = flag.String("org", "automata-devops-io", "organization to target in github")
)

func repoMan(w http.ResponseWriter, r *http.Request) {
	ghtoken := os.Getenv("GHTOKEN")
	whsecret := os.Getenv("WHSECRET")
	ctx := context.Background()
	// flag.Parse()
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghtoken},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	client := github.NewClient(tokenClient)

	payload, err := github.ValidatePayload(r, []byte(whsecret))
	if err != nil {
		log.Error("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#push
		// this is a commit push, do something with it
	case *github.PullRequestEvent:
		// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request
		// this is a pull request, do something with it
	case *github.WatchEvent:
		// https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#watch
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			log.Info("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	case *github.RepositoryEvent:
		//https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository
		// this is a repository event
		// this is where we manage the security settings
		if e.Action != nil && *e.Action == "created" {
			log.Info("new repository created. configuring security %s\n")
			opt := &github.RepositoryContentFileOptions{
				Message:   github.String("initial commit"),
				Content:   []byte(*github.String("# " + *e.Repo.Name)),
				Branch:    github.String("main"),
				Committer: &github.CommitAuthor{Name: github.String("Jeff Brimager"), Email: github.String("jbrimager@automata-devops.io")},
			}
			issue := &github.IssueRequest{
				Title:    github.String("New repo Created"),
				Body:     github.String("@sam1el this repo was created with the following rules applied\n - Require Pull Request Review\n - Requires 2 Approvers\n - Dismiss Stale Reviews\n - Require CodeOwner Review"),
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
			client.Repositories.CreateFile(ctx, *org, *e.Repo.Name, "README.md", opt)
			client.Repositories.UpdateBranchProtection(ctx, *org, *e.Repo.Name, "main", preq)
			client.Repositories.AddAdminEnforcement(ctx, *org, *e.Repo.Name, "main")
			client.Issues.Create(ctx, *org, *e.Repo.Name, issue)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
