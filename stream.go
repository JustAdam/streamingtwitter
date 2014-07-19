// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT

// Package streamingtwitter provides access to Twitter's streaming API.
// See https://dev.twitter.com/docs/api/streaming for more information.

package streamingtwitter

import (
	"encoding/json"
	"net"
	"net/url"
)

var (
	// Streams is a map of known Twitter Streaming API URLs.
	/* @todo implement fully */
	Streams = make(map[string]*TwitterAPIURL)
)

// https://dev.twitter.com/docs/api/streaming
func init() {
	// Public stream URLs - https://dev.twitter.com/docs/streaming-apis/streams/public
	Streams["Filter"] = &TwitterAPIURL{
		AccessMethod: "post",
		URL:          "https://stream.twitter.com/1.1/statuses/filter.json",
		Type:         "stream",
	}
	Streams["Firehose"] = &TwitterAPIURL{
		AccessMethod: "get",
		URL:          "https://stream.twitter.com/1.1/statuses/firehose.json",
		Type:         "stream",
	}
	Streams["Sample"] = &TwitterAPIURL{
		AccessMethod: "get",
		URL:          "https://stream.twitter.com/1.1/statuses/sample.json",
		Type:         "stream",
	}
}

// Stream creates a new Twitter API stream and sends received tweets on channel client.Tweets
/*
 args := &url.Values{}
 args.Add("track", "Norway")
 tweets := make(chan *streamingtwitter.TwitterStatus)
 go client.Stream(tweets, streamingtwitter.Streams["Filter"], args)
 for {
 	select {
		case status := <-tweets:
			fmt.Println(status)
		case err := <-client.Errors:
			fmt.Printf("ERROR: '%s'\n", err)
 		case <-client.Finished:
			return
		}
	}
*/
func (s *StreamClient) Stream(tweets chan<- *TwitterStatus, stream *TwitterAPIURL, formValues *url.Values) {
	resp, err := s.sendRequest(stream, formValues)
	if err != nil {
		s.Errors <- err
		return
	}
	defer func() {
		resp.Body.Close()
		s.Finished <- struct{}{}
	}()

	status := new(TwitterStatus)
	decoder := json.NewDecoder(resp.Body)
	for {
		// @todo handle these: https://dev.twitter.com/docs/streaming-apis/messages
		// @todo Handle stall_warnings if the option is set
		// @todo Handle fragmented JSON, (delimited)

		if err := decoder.Decode(&status); err != nil {
			if rerr, ok := err.(*net.OpError); ok {
				s.Errors <- rerr
				return
			} else if err.Error() == "EOF" {
				s.Errors <- err
				return
			} else if err.Error() == "unexpected EOF" {
				// Reconnection is left up to the client.
				s.Errors <- err
				return
			}
			s.Errors <- err
			continue
		}

		tweets <- status
	}
}
