// Sample Twitter client using the streaming twitter API.
package main

import (
	"fmt"
	"github.com/JustAdam/streamingtwitter"
	"log"
	"net/url"
)

var (
	// Location to the file containing your app's Twitter API token information
	tokenFile = "../tokens.json"
)

func main() {

	// Create new streaming API client
	client := streamingtwitter.NewClient()

	err := client.Authenticate(&streamingtwitter.ClientTokens{
		TokenFile: tokenFile,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Some keywords to track .. see the Twitter Streaming API documentation for more information
	args := &url.Values{}
	args.Add("track", "Norway")

	// Launch the stream
	tweets := make(chan *streamingtwitter.TwitterStatus)
	go client.Stream(tweets, streamingtwitter.Streams["Filter"], args)

	for {
		select {
		// Recieve tweets
		case status := <-tweets:
			fmt.Println(status)
			// Any errors that occured
		case err := <-client.Errors:
			fmt.Printf("ERROR: '%s'\n", err)
			// Stream has finished
		case <-client.Finished:
			return
		}
	}
}
