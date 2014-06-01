// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT

// Package streamingtwitter provides access to Twitter's streaming API.
// See https://dev.twitter.com/docs/api/streaming for more information.
package streamingtwitter

import (
	"encoding/json"
	"net/url"
)

var (
	// Streaming API URLs
	Streams = make(map[string]*TwitterApiUrl)
)

// https://dev.twitter.com/docs/api/streaming
func init() {
	// Public stream URLs - https://dev.twitter.com/docs/streaming-apis/streams/public
	Streams["Filter"] = &TwitterApiUrl{
		AccessMethod: "post",
		Url:          "https://stream.twitter.com/1.1/statuses/filter.json",
		Type:         "stream",
	}
	Streams["Firehose"] = &TwitterApiUrl{
		AccessMethod: "get",
		Url:          "https://stream.twitter.com/1.1/statuses/firehose.json",
		Type:         "stream",
	}
	Streams["Sample"] = &TwitterApiUrl{
		AccessMethod: "get",
		Url:          "https://stream.twitter.com/1.1/statuses/sample.json",
		Type:         "stream",
	}
}

// Create new Twitter stream.
//
// args := &url.Values{}
// args.Add("track", "Norway")
// go client.Stream(streamingtwitter.Streams["Filter"], args)
// for {
// 	select {
//		case status := <-client.Tweets:
//			fmt.Println(status)
//		case err := <-client.Errors:
//			fmt.Printf("ERROR: '%s'\n", err)
// 		case <-client.Finished:
//			return
//		}
//	}
func (s *StreamClient) Stream(stream *TwitterApiUrl, formValues *url.Values) {
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
			s.Errors <- err
			if err.Error() == "EOF" {
				return
			}
			continue
		}

		// Do we need to know which stream the tweet came from?
		s.Tweets <- status
	}
}
