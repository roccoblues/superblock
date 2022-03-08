package main

import (
	"fmt"
	"regexp"

	"github.com/dustin/go-humanize"
)

var tweetURLRegexp = regexp.MustCompile(`https?://twitter.com/\w+/status/(\d+)`)

func (app *application) block(tweetURL string) error {
	matches := tweetURLRegexp.FindStringSubmatch(tweetURL)
	if len(matches) != 2 {
		return fmt.Errorf("invalid tweet url '%s'", tweetURL)
	}
	tweetID := matches[1]

	user, err := app.client.Me(nil)
	if err != nil {
		return fmt.Errorf("failed to fetch user information: %w", err)
	}
	app.logger.Printf("Logged in as %s\n", user.Name)

	tweet, err := app.client.Tweet(tweetID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch tweet information: %w", err)
	}
	app.logger.Printf("Fetched tweet: %s\n", tweet.Text)

	author, err := app.client.User(tweet.AuthorID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch author information: %w", err)
	}

	blocked, err := app.client.Block(user.ID, author.ID)
	if err != nil {
		return fmt.Errorf("failed to block author: %w", err)
	}
	if blocked {
		app.logger.Printf("Blocked tweet author %s (%s) created %s\n", author.Name, author.ID, humanize.Time(author.CreatedAt))
	} else {
		app.logger.Printf("Author %s not blocked!\n", author.ID)
	}

	res := app.client.LikingUsers(tweet.ID, nil)
	var count int
	for res.NextPage() {
		users, err := res.User()
		if err != nil {
			return fmt.Errorf("failed to fetch liking users for tweet %s: %w", tweet.ID, err)
		}

		for _, u := range users {
			blocked, err := app.client.Block(user.ID, u.ID)
			if err != nil {
				return fmt.Errorf("failed to block user: %w", err)
			}
			if blocked {
				app.logger.Printf("Blocked user %s (%s) created %s\n", u.Name, u.ID, humanize.Time(u.CreatedAt))
			} else {
				app.logger.Printf("User %s not blocked!\n", u.ID)
			}
			count++
		}
	}

	app.logger.Printf("Blocked %d users who liked the tweet %s\n", count, tweetURL)

	return nil
}
