// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"bytes"
	"errors"
	"github.com/garyburd/go-oauth/oauth"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestUserLookupJsonDecode(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		cf, err := os.Open("test_data/user_lookup.json")
		if err != nil {
			t.Fatal("Unable to open test data file")
		}
		resp := &http.Response{
			Body: cf,
		}
		return resp, nil
	}

	testurl := &TwitterApiUrl{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	data := []TwitterUser{}
	go client.Rest(testurl, &url.Values{}, &data)

	select {
	case <-client.Finished:
		break
	case <-time.After(2 * time.Millisecond):
		t.Fatal("Data not received on Finished channel")
	}

	testData := []JsonTestData{
		{"Id", data[0].Id, "89409855"},
		{"Id", data[1].Id, "15439395"},
	}

	for _, d := range testData {
		if d.v != d.e {
			t.Errorf("%v: expecting %v, got %v", d.n, d.e, d.v)
		}
	}
}

// sendRequest error
func TestRestSendsRequestError(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		return &http.Response{}, errors.New("Test error")
	}

	testurl := &TwitterApiUrl{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	go client.Rest(testurl, &url.Values{}, &struct{}{})
	select {
	case err := <-client.Errors:
		if err.Error() != "Test error" {
			t.Errorf("Expecting error \"Test error\", got %v", err)
		}
		break
	case <-time.After(2 * time.Millisecond):
		t.Error("Error not received on Errors channel")
	}
}

// Decoding error
func TestRestSendsDecodingError(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("{invalid}")),
		}
		return resp, nil
	}

	testurl := &TwitterApiUrl{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	go client.Rest(testurl, &url.Values{}, &struct{}{})
	select {
	case err := <-client.Errors:
		if err.Error() != "invalid character 'i' looking for beginning of object key string" {
			t.Errorf("Expecting error \"invalid character 'i' looking for beginning of object key string\", got %v", err)
		}
		break
	case <-time.After(2 * time.Millisecond):
		t.Error("Error not received on Errors channel")
	}
}

func TestRestSendsFinishedNotification(t *testing.T) {
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("{\"x\":1}")),
		}
		return resp, nil
	}

	testurl := &TwitterApiUrl{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	go client.Rest(testurl, &url.Values{}, &struct{}{})
	select {
	case <-client.Finished:
		break
	case <-time.After(2 * time.Millisecond):
		t.Error("Data not received on Finished channel")
	}
}

func TestRestClosesResp(t *testing.T) {
	closedChannel := make(chan struct{})
	handler := func(*http.Client, *oauth.Credentials, string, url.Values) (*http.Response, error) {
		resp := &http.Response{
			Body: CloseCalled{
				bytes.NewBufferString("{\"x\":1}"),
				closedChannel,
			},
		}
		return resp, nil
	}

	testurl := &TwitterApiUrl{
		AccessMethod:  "custom",
		CustomHandler: handler,
	}

	client := NewClient()
	go client.Rest(testurl, &url.Values{}, &struct{}{})
	select {
	case <-closedChannel:
		break
	case <-time.After(2 * time.Millisecond):
		t.Error("Resp.body was not closed")
	}
}

type CloseCalled struct {
	io.Reader
	Y chan struct{}
}

func (c CloseCalled) Close() error {
	c.Y <- struct{}{}
	return nil
}
