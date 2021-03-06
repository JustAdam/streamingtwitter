// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"bytes"
	"errors"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

type JSONTestData struct {
	n string      // Variable name
	v interface{} // Variable value
	e interface{} // Expected
}

func TestTweetCreation(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		cf, err := os.Open("test_data/tweet.json")
		if err != nil {
			t.Fatal("Unable to open test data file")
		}
		resp := &http.Response{
			Body: cf,
		}
		return resp, nil
	}

	testurl := &TwitterAPIURL{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	status := new(TwitterStatus)
	tweets := make(chan *TwitterStatus)
	go client.Stream(tweets, testurl, &url.Values{})
	select {
	case status = <-tweets:
		break
	case <-time.After(2 * time.Millisecond):
		t.Fatal("Tweet data not receieved")
	}

	testData := []JSONTestData{
		{"ID", status.ID, "468728009768579073"},
		{"ReplyToStatusIDStr", status.ReplyToStatusIDStr, ""},
		{"ReplyToUserIDStr", status.ReplyToUserIDStr, ""},
		{"ReplyToUserScreenName", status.ReplyToUserScreenName, ""},
		{"CreatedAt", status.CreatedAt.T.String(), "2014-05-20 12:20:40 +0000 UTC"},
		{"Text", status.Text, "RT @IgorjI4J5: https://t.co/WM4MOusVww #RuinAToy Схема башни московского кремля"},
		// TwitterUser
		{"User.ID", status.User.ID, "2450810000"},
		{"User.Name", status.User.Name, "Афанасия Кудряшова"},
		{"User.ScreenName", status.User.ScreenName, "AfanasiyaI2E"},
		{"User.CreatedAt", status.User.CreatedAt.T.String(), "2014-04-18 04:12:25 +0000 UTC"},
		{"User.Location", status.User.Location, ""},
		{"User.URL", status.User.URL, ""},
		{"User.Description", status.User.Description, ""},
		{"User.Protected", status.User.Protected, false},
		{"User.FollowersCount", status.User.FollowersCount, uint32(7)},
		{"User.FriendsCount", status.User.FriendsCount, uint32(13)},
		{"User.ListedCount", status.User.ListedCount, uint32(0)},
		{"User.FavouritesCount", status.User.FavouritesCount, uint32(0)},
		{"User.StatusCount", status.User.StatusCount, uint32(139)},
		{"User.UtcOffset", status.User.UtcOffset, int32(0)},
		{"User.Timezone", status.User.Timezone, ""},
		{"User.GeoEnabled", status.User.GeoEnabled, false},
		{"User.Verified", status.User.Verified, false},
		{"User.Language", status.User.Language, "ru"},
		{"User.ContributorsEnabled", status.User.ContributorsEnabled, false},
		{"User.IsTranslator", status.User.IsTranslator, false},
		{"User.IsTranslationEnabled", status.User.IsTranslationEnabled, false},
		{"User.FollowRequestSent", status.User.FollowRequestSent, false},
		{"User.ProfileBackgroundColor", status.User.ProfileBackgroundColor, "C0DEED"},
		{"User.ProfileBackgroundImageURL", status.User.ProfileBackgroundImageURL, "http://abs.twimg.com/images/themes/theme1/bg.png"},
		{"User.ProfileBackgroundImageURLHttps", status.User.ProfileBackgroundImageURLHttps, "https://abs.twimg.com/images/themes/theme1/bg.png"},
		{"User.ProfileBackgroundTile", status.User.ProfileBackgroundTile, false},
		{"User.ProfileImageURL", status.User.ProfileImageURL, "http://abs.twimg.com/sticky/default_profile_images/default_profile_2_normal.png"},
		{"User.ProfileImageURLHttps", status.User.ProfileImageURLHttps, "https://abs.twimg.com/sticky/default_profile_images/default_profile_2_normal.png"},
		{"User.ProfileLinkColor", status.User.ProfileLinkColor, "0084B4"},
		{"User.ProfileSidebarBorderColor", status.User.ProfileSidebarBorderColor, "C0DEED"},
		{"User.ProfileSidebarFillColor", status.User.ProfileSidebarFillColor, "DDEEF6"},
		{"User.ProfileTextColor", status.User.ProfileTextColor, "333333"},
		{"User.ProfileUseBackgroundImage", status.User.ProfileUseBackgroundImage, true},
		{"User.DefaultProfile", status.User.DefaultProfile, true},
		{"User.DefaultProfileImage", status.User.DefaultProfileImage, true},
		{"Source", status.Source, "web"},
		{"Truncated", status.Truncated, false},
		{"Favorited", status.Favorited, false},
		{"Retweeted", status.Retweeted, false},
		{"RetweetedStatus[\"id_str\"]", status.RetweetedStatus["id_str"], "468660117211455488"},
		{"PossiblySensitive", status.PossiblySensitive, false},
		{"Language", status.Language, "bg"},
		{"RetweetCount", status.RetweetCount, uint32(0)},
		{"FavoriteCount", status.FavoriteCount, uint32(0)},
		// TwitterCoordinate
		{"Coordinates.Type", status.Coordinates.Type, "Point"},
		{"Coordinates.Coordinates[0]", status.Coordinates.Coordinates[0], -74.210251},
		{"Coordinates.Coordinates[1]", status.Coordinates.Coordinates[1], 40.422551},
		// TwitterPlace
		{"Place.ID", status.Place.ID, "27485069891a7938"},
		{"Place.URL", status.Place.URL, "https://api.twitter.com/1.1/geo/id/27485069891a7938.json"},
		{"Place.PlaceType", status.Place.PlaceType, "city"},
		{"Place.Name", status.Place.Name, "New York"},
		{"Place.FullName", status.Place.FullName, "New York, NY"},
		{"Place.CountryCode", status.Place.CountryCode, "US"},
		{"Place.Country", status.Place.Country, "United States"},
		{"Place.BoundingBox.Type", status.Place.BoundingBox.Type, "Polygon"},
		{"Place.BoundingBox.Coordinates[0].([]interface{})[0].([]interface{})[0]", status.Place.BoundingBox.Coordinates[0].([]interface{})[0].([]interface{})[0], -74.04725},
		{"Place.BoundingBox.Coordinates[0].([]interface{})[0].([]interface{})[2]", status.Place.BoundingBox.Coordinates[0].([]interface{})[0].([]interface{})[1], 40.541722},
		{"Place.BoundingBox.Coordinates[0].([]interface{})[2].([]interface{})[0]", status.Place.BoundingBox.Coordinates[0].([]interface{})[2].([]interface{})[0], -73.699793},
		{"Place.BoundingBox.Coordinates[0].([]interface{})[2].([]interface{})[1]", status.Place.BoundingBox.Coordinates[0].([]interface{})[2].([]interface{})[1], 40.91533},
		// TwitterEntity
		// TwitterHashTag
		{"Entities.Hashtags[0].Text", status.Entities.Hashtags[0].Text, "RuinAToy"},
		{"Entities.Hashtags[0].Indices[0]", status.Entities.Hashtags[0].Indices[0], uint(39)},
		{"Entities.Hashtags[0].Indices[1]", status.Entities.Hashtags[0].Indices[1], uint(48)},
		// TwitterMedia
		{"Entities.Media[0].ID", status.Entities.Media[0].ID, "468831177185710081"},
		{"Entities.Media[0].Type", status.Entities.Media[0].Type, "photo"},
		{"Entities.Media[0].URL", status.Entities.Media[0].URL, "http://t.co/Sqk7VYgixB"},
		{"Entities.Media[0].DisplayURL", status.Entities.Media[0].DisplayURL, "pic.twitter.com/Sqk7VYgixB"},
		{"Entities.Media[0].ExpandedURL", status.Entities.Media[0].ExpandedURL, "http://twitter.com/ShawnWiora/status/468831180297887744/photo/1"},
		{"Entities.Media[0].MediaURL", status.Entities.Media[0].MediaURL, "http://pbs.twimg.com/media/BoGfeL_IUAEcnfI.jpg"},
		{"Entities.Media[0].MediaURLHttps", status.Entities.Media[0].MediaURLHttps, "https://pbs.twimg.com/media/BoGfeL_IUAEcnfI.jpg"},
		{"Entities.Media[0].Sizes[\"medium\"].(map[string]interface{})[\"w\"]", status.Entities.Media[0].Sizes["medium"].(map[string]interface{})["w"], float64(600)},
		{"Entities.Media[0].Indices", status.Entities.Media[0].Indices[0], uint(113)},
		{"Entities.Media[0].Indices", status.Entities.Media[0].Indices[1], uint(135)},
		// TwitterUrl
		{"Entities.URLs[0].Url", status.Entities.URLs[0].URL, "https://t.co/WM4MOusVww"},
		{"Entities.URLs[0].DisplayURL", status.Entities.URLs[0].DisplayURL, "docs.google.com/document/d/1Uc\u2026"},
		{"Entities.URLs[0].ExpandedURL", status.Entities.URLs[0].ExpandedURL, "https://docs.google.com/document/d/1UcPtaqzHLCpbOlbrOBRAHuWoXsxoc2Ei4hxlNJMR_1I/edit?usp=d"},
		{"Entities.URLs[0].Indices", status.Entities.URLs[0].Indices[0], uint(15)},
		{"Entities.URLs[0].Indices", status.Entities.URLs[0].Indices[1], uint(38)},
		// TwitterMention
		{"Entities.UserMentions[0].ID", status.Entities.UserMentions[0].ID, "2459254292"},
		{"Entities.UserMentions[0].Name", status.Entities.UserMentions[0].Name, "Игорь Савин"},
		{"Entities.UserMentions[0].ScreenName", status.Entities.UserMentions[0].ScreenName, "IgorjI4J5"},
		{"Entities.UserMentions[0].Indices", status.Entities.UserMentions[0].Indices[0], uint(3)},
		{"Entities.UserMentions[0].Indices", status.Entities.UserMentions[0].Indices[1], uint(13)},
	}

	for _, d := range testData {
		if d.v != d.e {
			t.Errorf("%v: expecting %v, got %v", d.n, d.e, d.v)
		}
	}
}

func TestDefaultStreamVariablesExist(t *testing.T) {
	_, ok := Streams["Filter"]
	if ok != true {
		t.Error("Missing default stream: Filter")
	}
	_, ok = Streams["Firehose"]
	if ok != true {
		t.Error("Missing default stream: Firehose")
	}
	_, ok = Streams["Sample"]
	if ok != true {
		t.Error("Missing default stream: Sample")
	}
}

func TestStreamSendsRequestError(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		return &http.Response{}, errors.New("test error")
	}

	testurl := &TwitterAPIURL{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	tweets := make(chan *TwitterStatus)
	go client.Stream(tweets, testurl, &url.Values{})
	select {
	case err := <-client.Errors:
		if err.Error() != "test error" {
			t.Errorf("Expecting error \"Test error\", got %v", err)
		}
		break
	case <-time.After(2 * time.Millisecond):
		t.Error("Error not received on Errors channel")
	}
}

func TestStreamEOFClosesResp(t *testing.T) {
	closedChannel := make(chan struct{})
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: CloseCalled{
				bytes.NewBufferString("{\"x\": 1}"),
				closedChannel,
			},
		}
		return resp, nil
	}

	testurl := &TwitterAPIURL{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	tweets := make(chan *TwitterStatus)
	go client.Stream(tweets, testurl, &url.Values{})
	timeout := time.After(5 * time.Millisecond)
	errors := 0
	for {
		select {
		// Receive tweet
		case <-tweets:
			// Receive EOF
		case <-client.Errors:
			errors++
			if errors > 1 {
				t.Error("Stream was not closed after receiving EOF")
				return
			}
		case <-closedChannel:
			return
		case <-timeout:
			t.Error("Resp.body was not closed")
			return
		}
	}
}

func TestDecodingErrorContinues(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("{\"id_str\":\"1\",\"text\":\"text\",\"source\":\"text\"}\n\n{\"id_str\":\"1.1\",\"text\":234}\n\n{\"id_str\":\"2\",\"text\":\"text\"}\n\n")),
		}
		return resp, nil
	}

	testurl := &TwitterAPIURL{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	tweets := make(chan *TwitterStatus)
	go client.Stream(tweets, testurl, &url.Values{})
	timeout := time.After(5 * time.Millisecond)
	recTweets := 0
	for {
		select {
		case <-tweets:
			recTweets++
		case <-client.Errors:
		case <-client.Finished:
			if recTweets != 2 {
				t.Error("Decoding error did not continue")
			}
			return
		case <-timeout:
			t.Error("Decoding error timeout")
			return
		}
	}
}

// @todo fix: this test relies on resp.Body.Close() being called in Stream
func TestNetworkErrorReturns(t *testing.T) {
	closedChannel := make(chan struct{})
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: CloseCalled{
				bytes.NewBufferString("{\"x\": 1}"),
				closedChannel,
			},
		}
		return resp, &net.OpError{Err: errors.New("network error")}
	}

	testurl := &TwitterAPIURL{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	tweets := make(chan *TwitterStatus)
	go client.Stream(tweets, testurl, &url.Values{})
	select {
	case err := <-client.Errors:
		if _, ok := err.(*net.OpError); !ok {
			t.Error("Expecting error type &net.OpError")
		}
	case <-closedChannel:
		return
	case <-time.After(2 * time.Millisecond):
		t.Error("Error not received on Errors channel")
	}
}
