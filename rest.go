// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT

// Package streamingtwitter provides access to Twitter's streaming API.
// See https://dev.twitter.com/docs/api/streaming for more information.

package streamingtwitter

import (
	"encoding/json"
	"net/url"
)

// Rest sends a REST request to Twitter's REST API:  https://dev.twitter.com/docs/api/1.1
/*
 args := &url.Values{}
 args.Add("screen_name", "TwitterName")
 data := []streamingtwitter.TwitterUser{}
 url := &TwitterApiUrl{
  AccessMethod: "get",
  Url:          "https://api.twitter.com/1.1/users/lookup.json",
 }
 go client.Rest(&data, url, args)
 select {
 case err := <-client.Errors:
	log.Fatal(err)
 case <-client.Finished:
	fmt.Printf("%+v", data)
 }
*/
func (s *StreamClient) Rest(data interface{}, stream *TwitterAPIURL, formValues *url.Values) {
	resp, err := s.sendRequest(stream, formValues)
	if err != nil {
		s.Errors <- err
		return
	}
	defer func() {
		resp.Body.Close()
		s.Finished <- struct{}{}
	}()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		s.Errors <- err
		return
	}
}
