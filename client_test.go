// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"testing"
)

func TestAuthenticateMissingTokenDataError(t *testing.T) {
	client := NewClient()

	file := "test_data/tokens.json"
	err := client.Authenticate(&file)
	if err.Error() != "Missing App token" {
		t.Errorf("Expecting error \"Missing App token\", got %v", err)
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
