package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/google/go-github/github"
	"github.com/plutov/yearofcommits/icon"
	"golang.org/x/oauth2"
)

const dateFormat = "2016-01-02"

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetIcon(icon.Data)
	mQuit := systray.AddMenuItem("Quit", "")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		user := flag.String("u", "", "GitHub username")
		token := flag.String("t", "", "GitHub API token")
		flag.Parse()
		if *user == "" || *token == "" {
			fmt.Println("Usage: yearofcommits -u github-username -t github-api-token")
			systray.Quit()
			os.Exit(0)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *token},
		)
		tc := oauth2.NewClient(ctx, ts)
		c := github.NewClient(tc)

		// Check github stats every hour
		updateCounter(ctx, c, *user)
		for range time.Tick(time.Hour) {
			updateCounter(ctx, c, *user)
		}
	}()
}

func updateCounter(ctx context.Context, client *github.Client, user string) {
	repos, _, reposErr := client.Repositories.List(ctx, user, nil)
	if reposErr != nil {
		log.Printf("unable to get repos: %v", reposErr)
		return
	}

	today := time.Now()
	yearAgo := time.Now().AddDate(-1, 0, 0)

	commitsOpts := &github.CommitsListOptions{
		Author: user,
		Since:  yearAgo,
		Until:  today,
	}

	// date -> count of commits map
	dateCountMap := make(map[string]int)

	for _, repo := range repos {
		commits, _, commitsErr := client.Repositories.ListCommits(ctx, user, repo.GetName(), commitsOpts)
		if commitsErr != nil {
			log.Printf("unable to get commits of %s: %v", repo.GetName(), commitsErr)
			continue
		}

		for _, commit := range commits {
			dateCountMap[commit.Commit.GetAuthor().GetDate().Format(dateFormat)]++
		}
	}

	var daysInRow int
	day := time.Now()
	for day.After(yearAgo) {
		day = day.AddDate(0, 0, -1)
		if _, ok := dateCountMap[day.Format(dateFormat)]; ok {
			daysInRow++
		} else {
			// Stop loop
			break
		}
	}

	systray.SetTitle(fmt.Sprintf("%d", daysInRow))
}
