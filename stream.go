// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT

// Package streamingtwitter provides access to Twitter's streaming API.
// See https://dev.twitter.com/docs/api/streaming for more information.
package streamingtwitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	// File permissions for the token file.
	tokenFilePermission = os.FileMode(0600)
)

const (
	// Layout of Twitter's timestamp
	twitter_time_layout = "Mon Jan 02 15:04:05 Z0700 2006"
)

type StreamClient struct {
	oauthClient *oauth.Client
	token       *oauth.Credentials

	// Tweets from received from every open stream will be sent here
	Tweets chan *TwitterStatus
	// Any received errors are sent here (Embedded API errors are current not fully supported)
	Errors chan error
	// When a Stream call has finished, this channel will receive data
	Finished chan struct{}
}

// @todo add all the data defined by https://dev.twitter.com/docs/platform-objects/tweets
// - https://dev.twitter.com/docs/platform-objects/entities
// - https://dev.twitter.com/docs/platform-objects/places
type TwitterStatus struct {
	Id          string      `json:"id_str"`
	ReplyTo     string      `json:"in_reply_to_status_id_str"`
	ReplyToUser string      `json:"in_reply_to_user_id_str"`
	CreatedAt   TwitterTime `json:"created_at"`
	Text        string      `json:"text"`
	User        TwitterUser `json:"User"`
}

// Currently a seperate type for easier unmarshaling help
type TwitterTime struct {
	T time.Time
}

// @todo: https://dev.twitter.com/docs/platform-objects/users
type TwitterUser struct {
	Id         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func NewClient() (client *StreamClient) {
	client = new(StreamClient)
	client.oauthClient = &oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}

	client.Tweets = make(chan *TwitterStatus)
	client.Errors = make(chan error)
	client.Finished = make(chan struct{})
	return
}

// Authenicate the app and user with Twitter.
func (s *StreamClient) Authenticate(tokenFile *string) error {

	// Read in applications's token information. In json format:
	//	{
	//  "App":{
	//    "Token":"YOUR APP TOKEN HERE",
	//    "Secret":"APP SECRET HERE"
	//  }
	//}
	cf, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		return err
	}

	credentials := make(map[string]*oauth.Credentials)

	if err := json.Unmarshal(cf, &credentials); err != nil {
		return err
	}

	app, ok := credentials["App"]
	if ok != true {
		err = errors.New("Missing App token")
		return err
	}
	s.oauthClient.Credentials = *app

	// Check for token information from the user (they need to grant your app access
	// for feed access)
	token, ok := credentials["User"]
	if ok != true {

		tempCredentials, err := s.oauthClient.RequestTemporaryCredentials(http.DefaultClient, "oob", nil)
		if err != nil {
			return err
		}

		url := s.oauthClient.AuthorizationURL(tempCredentials, nil)
		fmt.Fprintf(os.Stdout, "Before we can continue ...\nGo to:\n\n\t%s\n\nAuthorize the application and enter in the verification code: ", url)

		var authCode string
		fmt.Scanln(&authCode)

		token, _, err := s.oauthClient.RequestToken(http.DefaultClient, tempCredentials, authCode)
		if err != nil {
			return err
		}

		credentials["User"] = token
		save, err := json.Marshal(credentials)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(*tokenFile, save, tokenFilePermission); err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "Auth token saved\n")
	}

	s.token = token

	return nil
}

// Create new Twitter stream
func (s *StreamClient) Stream(stream TwitterStream, formValues *url.Values) {

	var method func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error)
	if stream.AccessMethod == "post" {
		method = s.oauthClient.Post
	} else {
		method = s.oauthClient.Get
	}
	resp, err := method(http.DefaultClient, s.token, stream.Url, *formValues)
	if err != nil {
		s.Errors <- err
		return
	}
	defer func() {
		resp.Body.Close()
		s.Finished <- struct{}{}
	}()

	// https://dev.twitter.com/docs/streaming-api-response-codes
	switch resp.StatusCode {
	case 401:
		s.Errors <- errors.New("Incorrect usename or password.")
		// Delete User entry in tokens json file?
		return
	case 403:
		s.Errors <- errors.New("Access to resource is forbidden")
		return
	case 404:
		s.Errors <- errors.New("Resource does not exist.")
		return
	case 406:
		s.Errors <- errors.New("One or more required parameters are missing or are not suitable (see relevant stream API for more information).")
		return
	case 413:
		s.Errors <- errors.New("A parameter list is too long (contact Twitter for increased access).")
		return
	case 416:
		s.Errors <- errors.New("Range unacceptable.")
		return
	case 420:
		s.Errors <- errors.New("Rate limited.")
		return
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		// @todo handle these: https://dev.twitter.com/docs/streaming-apis/messages
		// @todo Handle stall_warnings if the option is set
		// @todo Handle errors if missing values are supplied
		// @todo Handle fragmented JSON, (delimited)

		status := new(TwitterStatus)
		if err := decoder.Decode(&status); err != nil {
			s.Errors <- err
			//if err.Error() == "Unexpected EOF" {
			//	return
			//}
			break
		}

		// Do we need to know which stream the tweet came from?
		s.Tweets <- status
	}
}

func (tt *TwitterTime) UnmarshalJSON(b []byte) (err error) {
	// Remove start and end quotes
	tt.T, err = time.Parse(twitter_time_layout, string(b[1:len(b)-1]))
	return
}
