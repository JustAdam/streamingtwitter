// Copyright 2014 JustAdam (adambell7@gmail.com).  All rights reserved.
// License: MIT
package streamingtwitter

var (
	// Streaming API URLs
	Streams = make(map[string]*TwitterStream)
)

type TwitterStream struct {
	// HTTP method which should be used to access the method
	AccessMethod string
	Url          string
}

// https://dev.twitter.com/docs/api/streaming
func init() {
	// Public stream URLs - https://dev.twitter.com/docs/streaming-apis/streams/public
	Streams["Filter"] = &TwitterStream{
		AccessMethod: "post",
		Url:          "https://stream.twitter.com/1.1/statuses/filter.json",
	}
	Streams["Firehose"] = &TwitterStream{
		AccessMethod: "get",
		Url:          "https://stream.twitter.com/1.1/statuses/firehose.json",
	}
	Streams["Sample"] = &TwitterStream{
		AccessMethod: "get",
		Url:          "https://stream.twitter.com/1.1/statuses/sample.json",
	}
}
