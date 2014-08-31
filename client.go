// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT

// Package streamingtwitter provides access to Twitter's streaming API.
// See https://dev.twitter.com/docs/api/streaming for more information.
package streamingtwitter

import (
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"net/http"
	"net/url"
	"time"
)

const (
	// Layout of Twitter's timestamp
	twitterTimeLayout = "Mon Jan 02 15:04:05 Z0700 2006"
)

// StreamClient provides a client to access to the Twitter API.  The client is unusable until
// it is authenticated with Twitter (call Authenticate()).
type StreamClient struct {
	oauthClient *oauth.Client
	token       *oauth.Credentials

	/* @todo Calling code should know which stream/request finishes or errors? */

	// Any received errors are sent here (Embedded API errors are currently not fully supported)
	Errors chan error
	// When a request has finished, this channel will receive data.
	Finished chan struct{}
}

// A TwitterAPIURL provides details on how to access Twitter API URLs.
type TwitterAPIURL struct {
	// HTTP method which should be used to access the method (currently only get, post & custom is supported)
	AccessMethod string
	// If setting AccessMethod to custom then you must provide your own client handler.  Otherwise all
	// requests go via the oauthClient.
	CustomHandler func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error)
	// An actual Twitter API resource URL.
	URL string
	// API type being accessed (stream or rest)
	Type string
}

// A TwitterError will be generated when there is a problem with the request or stream.
// JSON decoding errors are not changed.
type TwitterError struct {
	ID  int
	Msg string
}

// TwitterStatus represents a tweet with all supporting & available information.
type TwitterStatus struct {
	ID                    string                 `json:"id_str"`
	ReplyToStatusIDStr    string                 `json:"in_reply_to_status_id_str"`
	ReplyToUserIDStr      string                 `json:"in_reply_to_user_id_str"`
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

// TwitterTime provides a timestamp.  It is seperate for easier JSON unmarshaling help.
type TwitterTime struct {
	T time.Time
}

// TwitterUser represents a Twitter user with all supporting & available information.
type TwitterUser struct {
	ID                             string                 `json:"id_str"`
	Name                           string                 `json:"name"`
	ScreenName                     string                 `json:"screen_name"`
	CreatedAt                      TwitterTime            `json:"created_at"`
	Location                       string                 `json:"location"`
	URL                            string                 `json:"url"`
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
	ProfileBackgroundImageURL      string                 `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHttps string                 `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool                   `json:"profile_background_tile"`
	ProfileImageURL                string                 `json:"profile_image_url"`
	ProfileImageURLHttps           string                 `json:"profile_image_url_https"`
	ProfileLinkColor               string                 `json:"profile_link_color"`
	ProfileSidebarBorderColor      string                 `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string                 `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string                 `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool                   `json:"profile_use_background_image"`
	DefaultProfile                 bool                   `json:"default_profile"`
	DefaultProfileImage            bool                   `json:"default_profile_image"`
	Status                         map[string]interface{} `json:"status"`
}

// TwitterCoordinate is a Twitter platform object and stores coordinates.
type TwitterCoordinate struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
}

// TwitterPlace is a Twitter platform object for places.
type TwitterPlace struct {
	ID              string                 `json:"id"`
	URL             string                 `json:"url"`
	PlaceType       string                 `json:"place_type"`
	Name            string                 `json:"name"`
	FullName        string                 `json:"full_name"`
	CountryCode     string                 `json:"country_code"`
	Country         string                 `json:"country"`
	BoundingBox     TwitterCoordinate      `json:"bounding_box"`
	ContainedWithin map[string]interface{} `json:"contained_within"`
}

// TwitterEntity contains entity information associated to a tweet.
type TwitterEntity struct {
	Hashtags     []TweetHashTag     `json:"hashtags"`
	Media        []TweetMedia       `json:"media"`
	URLs         []TweetURL         `json:"urls"`
	UserMentions []TweetUserMention `json:"user_mentions"`
}

// TweetHashTag contains any hashtags that are found within the tweet.
type TweetHashTag struct {
	Text    string `json:"text"`
	Indices []uint `json:"indices"`
}

// TweetMedia contains any media types that are associated to the tweet.
type TweetMedia struct {
	ID             string                 `json:"id_str"`
	Type           string                 `json:"type"`
	URL            string                 `json:"url"`
	DisplayURL     string                 `json:"display_url"`
	ExpandedURL    string                 `json:"expanded_url"`
	MediaURL       string                 `json:"media_url"`
	MediaURLHttps  string                 `json:"media_url_https"`
	Sizes          map[string]interface{} `json:"sizes"` // https://dev.twitter.com/docs/platform-objects/entities#obj-sizes
	Indices        []uint                 `json:"indices"`
	SourceStatusID string                 `json:"source_id_status_str"`
}

// TweetURL contains any URLs that are found within the tweet.
type TweetURL struct {
	URL         string `json:"url"`
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []uint `json:"indices"`
}

// TweetUserMention containis any users who were mentioned within the tweet.
type TweetUserMention struct {
	ID         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Indices    []uint `json:"indices"`
}

// NewClient creates a new StreamClient for access to the Twitter API (both stream & rest).
func NewClient() (client *StreamClient) {
	client = new(StreamClient)
	client.oauthClient = &oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
	client.Errors = make(chan error)
	client.Finished = make(chan struct{})
	return
}

// Authenicate the app and user, with Twitter using the oauth client.
func (s *StreamClient) Authenticate(t Tokener) (err error) {
	s.token, err = t.Token(s.oauthClient)
	return
}

// Send a request to Twitter.
// Calling method is responsible for closing the connection.
func (s *StreamClient) sendRequest(stream *TwitterAPIURL, formValues *url.Values) (*http.Response, error) {
	var method func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error)
	if stream.AccessMethod == "custom" {
		method = stream.CustomHandler
	} else {
		if stream.AccessMethod == "post" {
			method = s.oauthClient.Post
		} else {
			method = s.oauthClient.Get
		}
	}

	resp, err := method(http.DefaultClient, s.token, stream.URL, *formValues)
	if err != nil {
		return nil, err
	}

	// https://dev.twitter.com/docs/streaming-api-response-codes
	switch resp.StatusCode {
	case 401:
		// Delete User entry in tokens json file?
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "Incorrect usename or password.",
		}
	case 403:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "Access to resource is forbidden",
		}
	case 404:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "Resource does not exist.",
		}
	case 406:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "One or more required parameters are missing or are not suitable (see relevant stream API for more information).",
		}
	case 413:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "A parameter list is too long (contact Twitter for increased access).",
		}
	case 416:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "Range unacceptable.",
		}
	case 420:
		return nil, &TwitterError{
			ID:  resp.StatusCode,
			Msg: "Rate limited.",
		}
	default:
		return resp, nil
	}
}

// Unmarshal a timestamp from Twitter
func (tt *TwitterTime) UnmarshalJSON(b []byte) (err error) {
	tt.T, err = time.Parse(twitterTimeLayout, string(b[1:len(b)-1]))
	return
}

func (e TwitterError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Msg, e.ID)
}
