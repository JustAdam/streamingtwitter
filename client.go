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

// @todo Calling code should know which stream/request finishes or errors
type StreamClient struct {
	oauthClient *oauth.Client
	token       *oauth.Credentials

	// Tweets received from every open stream will be sent here
	Tweets chan *TwitterStatus
	// Any received errors are sent here (Embedded API errors are current not fully supported)
	Errors chan error
	// When a call has finished, this channel will receive data
	Finished chan struct{}
}

type TwitterApiUrl struct {
	// HTTP method which should be used to access the method
	AccessMethod string
	Url          string
	// API type being accessed (stream or rest)
	Type string
}

type TwitterError struct {
	Id  int
	Msg string
}

type TwitterStatus struct {
	Id                    string                 `json:"id_str"`
	ReplyToStatusIdStr    string                 `json:"in_reply_to_status_id_str"`
	ReplyToUserIdStr      string                 `json:"in_reply_to_user_id_str"`
	ReplyToUserScreenName string                 `json:"in_reply_to_screen_name"`
	CreatedAt             TwitterTime            `json:"created_at"`
	Text                  string                 `json:"text"`
	User                  TwitterUser            `json:"User"`
	Source                string                 `json:"source"`
	Truncated             bool                   `json:"truncated"`
	Favorited             bool                   `json:"favorited"`
	Retweeted             bool                   `json:"retweeted"`
	RetweetedStatus       map[string]interface{} `json:"retweeted_status"`
	PossiblySensitive     bool                   `json:"possibly_sensitive"`
	Language              string                 `json:"lang"`
	RetweetCount          uint32                 `json:"retweet_count"`
	FavoriteCount         uint32                 `json:"favorite_count"`
	Coordinates           TwitterCoordinate      `json:"coordinates"`
	Place                 TwitterPlace           `json:"place"`
	Entities              TwitterEntity          `json:"entities"`
}

// Easier JSON unmarshaling help
type TwitterTime struct {
	T time.Time
}

type TwitterUser struct {
	Id                             string                 `json:"id_str"`
	Name                           string                 `json:"name"`
	ScreenName                     string                 `json:"screen_name"`
	CreatedAt                      TwitterTime            `json:"created_at"`
	Location                       string                 `json:"location"`
	Url                            string                 `json:"url"`
	Description                    string                 `json:"description"`
	Protected                      bool                   `json:"protected"`
	FollowersCount                 uint32                 `json:"followers_count"`
	FriendsCount                   uint32                 `json:"friends_count"`
	ListedCount                    uint32                 `json:"listed_count"`
	FavouritesCount                uint32                 `json:"favourites_count"`
	StatusCount                    uint32                 `json:"statuses_count"`
	UtcOffset                      int32                  `json:"utc_offset"`
	Timezone                       string                 `json:"time_zone"`
	GeoEnabled                     bool                   `json:"geo_enabled"`
	Verified                       bool                   `json:"verified"`
	Language                       string                 `json:"lang"`
	ContributorsEnabled            bool                   `json:"contributors_enabled"`
	IsTranslator                   bool                   `json:"is_translator"`
	IsTranslationEnabled           bool                   `json:"is_translation_enabled"`
	FollowRequestSent              bool                   `json:"follow_request_sent"`
	ProfileBackgroundColor         string                 `json:"profile_background_color"`
	ProfileBackgroundImageUrl      string                 `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string                 `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool                   `json:"profile_background_tile"`
	ProfileImageUrl                string                 `json:"profile_image_url"`
	ProfileImageUrlHttps           string                 `json:"profile_image_url_https"`
	ProfileLinkColor               string                 `json:"profile_link_color"`
	ProfileSidebarBorderColor      string                 `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string                 `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string                 `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool                   `json:"profile_use_background_image"`
	DefaultProfile                 bool                   `json:"default_profile"`
	DefaultProfileImage            bool                   `json:"default_profile_image"`
	Status                         map[string]interface{} `json:"status"`
}

type TwitterCoordinate struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
}

type TwitterPlace struct {
	Id              string                 `json:"id"`
	Url             string                 `json:"url"`
	PlaceType       string                 `json:"place_type"`
	Name            string                 `json:"name"`
	FullName        string                 `json:"full_name"`
	CountryCode     string                 `json:"country_code"`
	Country         string                 `json:"country"`
	BoundingBox     TwitterCoordinate      `json:"bounding_box"`
	ContainedWithin map[string]interface{} `json:"contained_within"`
}

type TwitterEntity struct {
	Hashtags     []TweetHashTag     `json:"hashtags"`
	Media        []TweetMedia       `json:"media"`
	Urls         []TweetUrl         `json:"urls"`
	UserMentions []TweetUserMention `json:"user_mentions"`
}

type TweetHashTag struct {
	Text    string `json:"text"`
	Indices []uint `json:"indices"`
}

type TweetMedia struct {
	Id             string                 `json:"id_str"`
	Type           string                 `json:"type"`
	Url            string                 `json:"url"`
	DisplayUrl     string                 `json:"display_url"`
	ExpandedUrl    string                 `json:"expanded_url"`
	MediaUrl       string                 `json:"media_url"`
	MediaUrlHttps  string                 `json:"media_url_https"`
	Sizes          map[string]interface{} `json:"sizes"` // https://dev.twitter.com/docs/platform-objects/entities#obj-sizes
	Indices        []uint                 `json:"indices"`
	SourceStatusId string                 `json:"source_id_status_str"`
}

type TweetUrl struct {
	Url         string `json:"url"`
	DisplayUrl  string `json:"display_url"`
	ExpandedUrl string `json:"expanded_url"`
	Indices     []uint `json:"indices"`
}

type TweetUserMention struct {
	Id         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Indices    []uint `json:"indices"`
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
	//  {
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
	if s.oauthClient.Credentials.Token == "" || s.oauthClient.Credentials.Secret == "" {
		return errors.New("Missing app's Token or Secret")
	}
	s.oauthClient.Credentials = *app

	// Check for token information from the user (they need to grant your app access for feed access)
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

// Calling method is responsible for closing the connection.
func (s *StreamClient) sendRequest(stream *TwitterApiUrl, formValues *url.Values) (*http.Response, error) {
	var method func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error)
	if stream.AccessMethod == "post" {
		method = s.oauthClient.Post
	} else {
		method = s.oauthClient.Get
	}
	resp, err := method(http.DefaultClient, s.token, stream.Url, *formValues)
	if err != nil {
		return nil, err
	}

	// https://dev.twitter.com/docs/streaming-api-response-codes
	switch resp.StatusCode {
	case 401:
		// Delete User entry in tokens json file?
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "Incorrect usename or password.",
		}
	case 403:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "Access to resource is forbidden",
		}
	case 404:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "Resource does not exist.",
		}
	case 406:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "One or more required parameters are missing or are not suitable (see relevant stream API for more information).",
		}
	case 413:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "A parameter list is too long (contact Twitter for increased access).",
		}
	case 416:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "Range unacceptable.",
		}
	case 420:
		return nil, &TwitterError{
			Id:  resp.StatusCode,
			Msg: "Rate limited.",
		}
	default:
		return resp, nil
	}
}

func (tt *TwitterTime) UnmarshalJSON(b []byte) (err error) {
	// Remove start and end quotes
	tt.T, err = time.Parse(twitter_time_layout, string(b[1:len(b)-1]))
	return
}

func (e TwitterError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Msg, e.Id)
}
