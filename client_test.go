// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"github.com/garyburd/go-oauth/oauth"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestAuthenticateMissingAppDataError(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens.json"
	err := client.Authenticate(&file)
	if err.Error() != "missing App token" {
		t.Errorf("Expecting error \"Missing App token\", got %v", err)
	}
}

func TestAuthenticateMissingAppTokenSecretError(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens_empty.json"
	err := client.Authenticate(&file)
	if err.Error() != "missing app's Token or Secret" {
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
		ID:  101,
		Msg: "error message",
	}

	if err.Error() != "error message (101)" {
		t.Errorf("Expecting \"error message (101)\", got %v", err)
	}
}

func TestSendResponseErrorOutput(t *testing.T) {
	client := NewClient()
	errorCodes := []int{401, 403, 404, 406, 413, 416, 420}

	for _, v := range errorCodes {
		handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: v,
			}
			return resp, nil
		}

		testurl := &TwitterAPIURL{
			AccessMethod:  "custom",
			CustomHandler: handler,
		}

		_, err := client.sendRequest(testurl, &url.Values{})

		if rerr, ok := err.(*TwitterError); !ok {
			t.Errorf("Expecting TwitterError, got %v", reflect.TypeOf(err))
		} else if rerr.ID != v {
			t.Errorf("Expecting error ID %v, got %v", v, rerr.ID)
		}
	}
}
