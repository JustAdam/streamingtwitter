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

func TestTwitterErrorOutput(t *testing.T) {
	err := &TwitterError{
		ID:  101,
		Msg: "error message",
	}

	if err.Error() != "error message (101)" {
		t.Errorf("Expecting \"error message (101)\", got %v", err)
	}
}

func TestAuthenticateSetsAccessToken(t *testing.T) {
	client := NewClient()

	_ = client.Authenticate(&ClientTokens{
		TokenFile: "test_data/tokens-empty.json",
		App: &oauth.Credentials{
			Token:  "app-token",
			Secret: "app-secret",
		},
		User: &oauth.Credentials{
			Token:  "user-token",
			Secret: "user-secret",
		},
	})

	if client.token.Token != "user-token" && client.token.Secret != "user-secret" {
		t.Errorf("Expecting client.token (.Token = user-token) (.Secret = user-secret), got (.Token = %v) (.Secret = %v)", client.token.Token, client.token.Secret)
	}
}

func TestAuthenticateError(t *testing.T) {
	client := NewClient()

	err := client.Authenticate(&ClientTokens{
		TokenFile: "test_data/tokens-empty.json",
		User: &oauth.Credentials{
			Token:  "user-token",
			Secret: "user-secret",
		},
	})

	if _, ok := err.(*ClientTokensError); !ok {
		t.Errorf("Expecting ClientTokensError got %v", reflect.TypeOf(err))
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
