package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/google/go-github/github"
	"github.com/plutov/yearofcommits/icon"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

// Config Go representation
type Config struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

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
		ymlFile, err := ioutil.ReadFile("config.yml")
		if err != nil {
			log.Fatalf("unable to open config.yml: %v", err)
		}

		cfg := new(Config)
		err = yaml.Unmarshal([]byte(ymlFile), &cfg)
		if err != nil {
			log.Fatalf("unable to parse config: %v", err)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		c := github.NewClient(tc)

		// Check github stats every hour
		updateCounter(ctx, c, cfg)
		for range time.Tick(time.Hour) {
			updateCounter(ctx, c, cfg)
		}
	}()
}

func updateCounter(ctx context.Context, client *github.Client, cfg *Config) {
	repos, _, reposErr := client.Repositories.List(ctx, cfg.User, nil)
	if reposErr != nil {
		log.Printf("unable to get repos: %v", reposErr)
		return
	}

	today := time.Now()
	yearAgo := time.Now().AddDate(-1, 0, 0)

	commitsOpts := &github.CommitsListOptions{
		Since: yearAgo,
		Until: today,
	}

	// date -> count of commits map
	dateCountMap := make(map[string]int)

	for _, repo := range repos {
		commits, _, commitsErr := client.Repositories.ListCommits(ctx, cfg.User, repo.GetName(), commitsOpts)
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
