// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"testing"
)

func TestAuthenticateMissingAppDataError(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens.json"
	err := client.Authenticate(&file)
	if err.Error() != "Missing App token" {
		t.Errorf("Expecting error \"Missing App token\", got %v", err)
	}
}

func TestAuthenticateMissingAppTokenSecretError(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens_empty.json"
	err := client.Authenticate(&file)
	if err.Error() != "Missing app's Token or Secret" {
		t.Errorf("Expecting error \"Missing app's Token or Secret\", got %v", err)
	}
}

func TestAuthenticateAccessTokenIsSetInFile(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens_full.json"
	client.Authenticate(&file)
	if client.token.Token != "user-token" || client.token.Secret != "user-secret" {
		t.Errorf("Client access token not set.")
	}
}

func TestTwitterErrorOutput(t *testing.T) {
	err := &TwitterError{
		Id:  101,
		Msg: "Error message",
	}

	if err.Error() != "Error message (101)" {
		t.Errorf("Expecting \"Error message (101)\", got %v", err)
	}
}
