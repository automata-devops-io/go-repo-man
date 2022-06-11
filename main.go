package main

import (
	"context"
	"flag"
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
	http.HandleFunc("/", repoList)
	http.HandleFunc("/webhook", repoMan)
	log.Fatal(http.ListenAndServe(port, nil))
}

var (
	org = flag.String("org", "automata-devops-io", "organization to target in github")
)

func repoList(w http.ResponseWriter, r *http.Request) {
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	client := github.NewClient(tokenClient)
	repoOpt := &github.RepositoryListByOrgOptions{Type: "all"}

	repoList, _, err := client.Repositories.ListByOrg(context, *org, repoOpt)
	for _, repo := range repoList {
		log.Printf("[DEBUG] Repo %s: %s\n", *repo.Owner.Login, *repo.Name)
	}
	if err != nil {
		log.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}
}

func repoMan(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	flag.Parse()
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "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"},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	client := github.NewClient(tokenClient)

	payload, err := github.ValidatePayload(r, []byte("my-secret-key"))
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
			fileContent := []byte("This is the content of my file\nand the 2nd line of it")
			opts := &github.RepositoryContentFileOptions{
				Message:   github.String("Initial commit"),
				Content:   fileContent,
				Branch:    github.String("main"),
				Committer: &github.CommitAuthor{Name: github.String("Jeff Brimager"), Email: github.String("jbrimager@gmail.com")},
			}
			preq := &github.ProtectionRequest{
				EnforceAdmins: true,
				Restrictions:  nil,
				RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
					RequiredApprovingReviewCount: 2,
				},
				RequiredStatusChecks: &github.RequiredStatusChecks{
					Strict: true,
				},
			}
			client.Repositories.CreateFile(ctx, *e.Org.Name, *e.Repo.Name, "README.md", opts)
			client.Repositories.UpdateBranchProtection(ctx, *e.Org.Name, *e.Repo.Name, "main", preq)
			client.Repositories.AddAdminEnforcement(ctx, *e.Org.Name, *e.Repo.Name, "main")
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
