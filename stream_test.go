// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"testing"
)

func TestTwitterTimeUnmarshal(t *testing.T) {
	time := &TwitterTime{}
	time.UnmarshalJSON([]byte("\"Sat Sep 04 16:10:54 +0000 2010\""))
	if time.T.String() != "2010-09-04 16:10:54 +0000 UTC" {
		t.Errorf("Expecting '2010-09-04 16:10:54 +0000 UTC', got %v", time.T)
	}
}
