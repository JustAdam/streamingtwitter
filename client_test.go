// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"testing"
)

func TestTwitterErrorOutput(t *testing.T) {
	err := &TwitterError{
		Id:  101,
		Msg: "Error message",
	}

	if err.Error() != "Error message (101)" {
		t.Errorf("Expecting %v, got %v", err, "Error message (101)")
	}
}
