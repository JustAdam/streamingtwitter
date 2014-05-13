// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"testing"
)

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
