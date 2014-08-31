// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
//
// Sample Twitter client using the streaming twitter API.
//
// Displaying the tweets according to the terminal's size is not yet done,
// instead it uses a hard limit.
package main

import (
	"flag"
	"fmt"
	"github.com/JustAdam/streamingtwitter"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// Number of tweets to display
	tweetListLimit = 10
	// In seconds (float64)
	tweetHighlightAge = 3.0
	// Date/time display (type: go layout)
	tweetDateLayout = "15:04:05"
	// Timezone settings
	timezone, _ = time.LoadLocation("Europe/Oslo")

	// Flag parsing variables
	tokenFile     string
	stream        string
	followUsers   string
	location      string
	trackKeywords string

	wg sync.WaitGroup
)

func init() {
	flag.StringVar(&tokenFile, "config", "../tokens.json", "Token storage file location")
	flag.StringVar(&stream, "stream", "Filter", "Type of stream to open: <Filter>, <Firehose>, <Sample>")
	flag.StringVar(&followUsers, "follow", "", "Twitter users to track seperated by commas")
	flag.StringVar(&trackKeywords, "track", "", "Keywords to track seperated by commas")
	flag.StringVar(&location, "location", "", "Longitude and latitude keypairs to track seperated by commas")
}

func main() {
	// Clear terminal screen and move cursor to top left position.
	var clearScreen = func() {
		fmt.Fprintf(os.Stdout, "\033[2J\033[H")
	}

	clearScreen()

	// Splash screen ...
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		i := 1
		dot := "."
		fmt.Fprintf(os.Stdout, "Loading ")
		for _ = range ticker.C {
			if i%4 == 0 {
				fmt.Fprintf(os.Stdout, "\033[3D   \033[3D")
			} else {
				fmt.Fprintf(os.Stdout, "%v", dot)
			}
			i++
		}
		clearScreen()
	}()

	flag.Parse()

	_, ok := streamingtwitter.Streams[stream]
	if ok == false {
		fmt.Fprintf(os.Stderr, "Usage of %v\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	// Helper function to create a streaming client only when necessary
	var client *streamingtwitter.StreamClient
	var createClient = func() {
		// Create new streaming API client
		client = streamingtwitter.NewClient()

		err := client.Authenticate(&streamingtwitter.ClientTokens{
			TokenFile: tokenFile,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Define arguments to pass to the stream.  (Only required Filter stream options are supported currently)
	// https://dev.twitter.com/docs/streaming-apis/parameters
	args := &url.Values{}
	if stream == "Filter" {
		// At least one of: follow,locations or track must be specified to use the filter stream
		if followUsers == "" && trackKeywords == "" && location == "" {
			fmt.Fprintf(os.Stderr, "Either; -follow, -track, or -location must be specified when using the Filter stream.\n")
			flag.PrintDefaults()
			return
		}

		createClient()

		if followUsers != "" {
			// To track a user, a user ID (and not screen name) must be sent to the stream.
			// To get a user ID we need to query the relevant Twitter REST API.
			users := &url.Values{}
			users.Add("screen_name", followUsers)

			userLookup := &streamingtwitter.TwitterAPIURL{
				AccessMethod: "get",
				URL:          "https://api.twitter.com/1.1/users/lookup.json",
			}

			data := []streamingtwitter.TwitterUser{}
			go client.Rest(&data, userLookup, users)

			select {
			case err := <-client.Errors:
				ticker.Stop()
				clearScreen()
				if err.(*streamingtwitter.TwitterError).ID == 404 {
					log.Fatalf("User %v doesn't exist", followUsers)
				} else {
					log.Fatal(err)
				}
			case <-client.Finished:
				break
			}

			ids := []string{}
			for _, o := range data {
				ids = append(ids, o.ID)
			}

			args.Add("follow", strings.Join(ids, ","))
		}
		if trackKeywords != "" {
			args.Add("track", trackKeywords)
		}
		if location != "" {
			args.Add("locations", location)
		}
	} else {
		createClient()
	}

	wg.Add(1)
	tweets := make(chan *streamingtwitter.TwitterStatus)
	go client.Stream(tweets, streamingtwitter.Streams[stream], args)

	// Wait for all streams to finish and then provide notification
	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	// Store last X number of tweets for outputting
	tweetList := make([]streamingtwitter.TwitterStatus, tweetListLimit)

	// Push item to the top of the tweetList array and remove the last element if > listsize
	var addToFront = func(v streamingtwitter.TwitterStatus, l []streamingtwitter.TwitterStatus) (t []streamingtwitter.TwitterStatus) {
		t = make([]streamingtwitter.TwitterStatus, cap(l))
		t[0] = v
		for i := 1; i < cap(l); i++ {
			t[i] = l[i-1]
		}
		return
	}

	// Whether or not any tweets have been matched by the age timing criteria
	highlightMatched := false

	// Set the background and foreground colour of tweets matching age timing criteria
	var displayTweets = func(bg int, fg int) {
		clearScreen()

		for k, s := range tweetList {
			// Reset all display attributes
			fmt.Fprintf(os.Stdout, "\033[0m")

			// Default colour for tweet display
			// Clear all terminal formatting
			colourCode := "0"

			// Perphaps the output looks a little more pretty with a blank line at the top!
			if k == 0 {
				fmt.Fprintf(os.Stdout, "\n")
			}

			// Highlight all tweets that are less that X seconds old
			if time.Now().Sub(s.CreatedAt.T).Seconds() <= tweetHighlightAge {
				fmt.Fprintf(os.Stdout, "\033[%vm\033[%vm", bg, fg)
				colourCode = "22"
				highlightMatched = true
			}

			if s.Text != "" {
				fmt.Fprintf(os.Stdout, "\033[1m%v @%v\033[%vm - %v\n> %v\n\n", s.User.Name, s.User.ScreenName, colourCode, s.CreatedAt.T.In(timezone).Format(tweetDateLayout), s.Text)
			}
		}
		// Move cursor to home position (upper left corner)
		fmt.Fprintf(os.Stdout, " \033[H")
	}

	// Whether or not we have received data from a stream
	loaded := false

	// Business end
	for {
		select {
		case status := <-tweets:
			if loaded == false {
				loaded = true
				ticker.Stop()
			}

			tweetList = addToFront(*status, tweetList)

			// Draw tweets
			displayTweets(41, 37)

			// Redraw the tweets if we haven't received one within X seconds (remove highlighting)
		case <-time.After(time.Duration(tweetHighlightAge) * time.Second):
			// Only redraw if tweets have matched the age timing criteria in default colours
			if highlightMatched == true {
				// Perhaps we should only redraw the affected tweets .. :)
				displayTweets(49, 39)
				highlightMatched = false
			}
		case err := <-client.Errors:
			fmt.Fprintf(os.Stderr, "ERROR: '%s'\n", err)
		case <-client.Finished:
			// Notify the waitgroup that a stream has finished
			wg.Done()
		case <-done:
			// All streams are finished so we can exit
			return
		}
	}
}
