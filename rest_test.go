// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

import (
	"encoding/json"
	"os"
	"testing"
)

func GetUserLookupData(t *testing.T) (testData []JsonTestData) {
	cf, err := os.Open("test_data/user_lookup.json")
	if err != nil {
		t.Fatal("Unable to open test data file")
	}

	data := []TwitterUser{}
	err = json.NewDecoder(cf).Decode(&data)
	if err != nil {
		t.Errorf("Decoding into []TwitterUser failed, %v", err)
	}

	testData = []JsonTestData{
		{"Id", data[0].Id, "89409855"},
		{"Id", data[1].Id, "15439395"},
	}
	return
}

func TestUserLookupJsonDecode(t *testing.T) {
	testData := GetUserLookupData(t)
	for _, d := range testData {
		if d.v != d.e {
			t.Errorf("%v: expecting %v, got %v", d.n, d.e, d.v)
		}
	}
}
