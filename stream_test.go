// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"encoding/json"
	"os"
	"testing"
)

type JsonTestData struct {
	n string      // Variable name
	v interface{} // Variable value
	e interface{} // Expected
}

func GetTweetData(t *testing.T) (status *TwitterStatus, testData []JsonTestData) {

	cf, err := os.Open("test_data/tweet.json")
	if err != nil {
		t.Fatal("Unable to open test data file")
	}

	decoder := json.NewDecoder(cf)
	status = new(TwitterStatus)
	if err := decoder.Decode(&status); err != nil {
		t.Errorf("Unmarshaing into TwitterStatus failed, %v", err)
	}

	testData = []JsonTestData{
		{"Id", status.Id, "468728009768579073"},
		{"ReplyToStatusIdStr", status.ReplyToStatusIdStr, ""},
		{"ReplyToUserIdStr", status.ReplyToUserIdStr, ""},
		{"ReplyToUserScreenName", status.ReplyToUserScreenName, ""},
		{"CreatedAt", status.CreatedAt.T.String(), "2014-05-20 12:20:40 +0000 UTC"},
		{"Text", status.Text, "RT @IgorjI4J5: https://t.co/WM4MOusVww #RuinAToy Схема башни московского кремля"},
		// TwitterUser
		{"User.Id", status.User.Id, "2450810000"},
		{"User.Name", status.User.Name, "Афанасия Кудряшова"},
		{"User.ScreenName", status.User.ScreenName, "AfanasiyaI2E"},
		{"User.CreatedAt", status.User.CreatedAt.T.String(), "2014-04-18 04:12:25 +0000 UTC"},
		{"User.Location", status.User.Location, ""},
		{"User.Url", status.User.Url, ""},
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
		{"User.ProfileBackgroundImageUrl", status.User.ProfileBackgroundImageUrl, "http://abs.twimg.com/images/themes/theme1/bg.png"},
		{"User.ProfileBackgroundImageUrlHttps", status.User.ProfileBackgroundImageUrlHttps, "https://abs.twimg.com/images/themes/theme1/bg.png"},
		{"User.ProfileBackgroundTile", status.User.ProfileBackgroundTile, false},
		{"User.ProfileImageUrl", status.User.ProfileImageUrl, "http://abs.twimg.com/sticky/default_profile_images/default_profile_2_normal.png"},
		{"User.ProfileImageUrlHttps", status.User.ProfileImageUrlHttps, "https://abs.twimg.com/sticky/default_profile_images/default_profile_2_normal.png"},
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
		{"Place.Id", status.Place.Id, "27485069891a7938"},
		{"Place.Url", status.Place.Url, "https://api.twitter.com/1.1/geo/id/27485069891a7938.json"},
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
		{"Entities.Media[0].Id", status.Entities.Media[0].Id, "468831177185710081"},
		{"Entities.Media[0].Type", status.Entities.Media[0].Type, "photo"},
		{"Entities.Media[0].Url", status.Entities.Media[0].Url, "http://t.co/Sqk7VYgixB"},
		{"Entities.Media[0].DisplayUrl", status.Entities.Media[0].DisplayUrl, "pic.twitter.com/Sqk7VYgixB"},
		{"Entities.Media[0].ExpandedUrl", status.Entities.Media[0].ExpandedUrl, "http://twitter.com/ShawnWiora/status/468831180297887744/photo/1"},
		{"Entities.Media[0].MediaUrl", status.Entities.Media[0].MediaUrl, "http://pbs.twimg.com/media/BoGfeL_IUAEcnfI.jpg"},
		{"Entities.Media[0].MediaUrlHttps", status.Entities.Media[0].MediaUrlHttps, "https://pbs.twimg.com/media/BoGfeL_IUAEcnfI.jpg"},
		{"Entities.Media[0].Sizes[\"medium\"].(map[string]interface{})[\"w\"]", status.Entities.Media[0].Sizes["medium"].(map[string]interface{})["w"], float64(600)},
		{"Entities.Media[0].Indices", status.Entities.Media[0].Indices[0], uint(113)},
		{"Entities.Media[0].Indices", status.Entities.Media[0].Indices[1], uint(135)},
		// TwitterUrl
		{"Entities.Urls[0].Url", status.Entities.Urls[0].Url, "https://t.co/WM4MOusVww"},
		{"Entities.Urls[0].DisplayUrl", status.Entities.Urls[0].DisplayUrl, "docs.google.com/document/d/1Uc\u2026"},
		{"Entities.Urls[0].ExpandedUrl", status.Entities.Urls[0].ExpandedUrl, "https://docs.google.com/document/d/1UcPtaqzHLCpbOlbrOBRAHuWoXsxoc2Ei4hxlNJMR_1I/edit?usp=d"},
		{"Entities.Urls[0].Indices", status.Entities.Urls[0].Indices[0], uint(15)},
		{"Entities.Urls[0].Indices", status.Entities.Urls[0].Indices[1], uint(38)},
		// TwitterMention
		{"Entities.UserMentions[0].Id", status.Entities.UserMentions[0].Id, "2459254292"},
		{"Entities.UserMentions[0].Name", status.Entities.UserMentions[0].Name, "Игорь Савин"},
		{"Entities.UserMentions[0].ScreenName", status.Entities.UserMentions[0].ScreenName, "IgorjI4J5"},
		{"Entities.UserMentions[0].Indices", status.Entities.UserMentions[0].Indices[0], uint(3)},
		{"Entities.UserMentions[0].Indices", status.Entities.UserMentions[0].Indices[1], uint(13)},
	}
	return
}

func TestTweetCreation(t *testing.T) {
	_, testData := GetTweetData(t)
	for _, d := range testData {
		if d.v != d.e {
			t.Errorf("%v: expecting %v, got %v", d.n, d.e, d.v)
		}
	}
}
