// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"github.com/garyburd/go-oauth/oauth"
	"testing"
)

func TestTokenAbsentTokenFileStringError(t *testing.T) {
	tokens := &ClientTokens{}
	_, err := tokens.Token(&oauth.Client{})
	if err.Error() != "no token file supplied" {
		t.Errorf("Expecting error \"no token file supplied\", got %v", err)
	}
}

func TestTokenMissingAppDataError(t *testing.T) {
	tokens := &ClientTokens{
		TokenFile: "test_data/tokens-empty.json",
	}
	_, err := tokens.Token(&oauth.Client{})
	if err.Error() != "missing \"App\" token" {
		t.Errorf("Expecting error \"missing \"App\" token\", got %v", err)
	}
}

func TestTokenMissingAppTokenSecretError(t *testing.T) {
	tokens := &ClientTokens{
		TokenFile: "test_data/tokens-empty.json",
		App:       &oauth.Credentials{},
		User:      &oauth.Credentials{},
	}
	_, err := tokens.Token(&oauth.Client{})
	if err.Error() != "missing app's Token or Secret" {
		t.Errorf("Expecting error \"Missing app's Token or Secret\", got %v", err)
	}
}

func TestTokenAccessTokenIsSetInFile(t *testing.T) {
	tokens := &ClientTokens{
		TokenFile: "test_data/tokens-empty.json",
		App: &oauth.Credentials{
			Token:  "app-token",
			Secret: "app-secret",
		},
		User: &oauth.Credentials{
			Token:  "user-token",
			Secret: "user-secret",
		},
	}
	token, _ := tokens.Token(&oauth.Client{})
	if token.Token != "user-token" || token.Secret != "user-secret" {
		t.Errorf("Client access token not set.")
	}
}

func TestTokenErrorOutput(t *testing.T) {
	err := &ClientTokensError{
		Msg: "error message",
	}

	if err.Error() != "error message" {
		t.Errorf("Expecting \"error message\", got %v", err)
	}
}
