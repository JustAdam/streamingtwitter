Go: Twitter Streaming and REST API package
==========================================

[![Build Status](https://travis-ci.org/JustAdam/streamingtwitter.svg?branch=master)](https://travis-ci.org/JustAdam/streamingtwitter) [![Coverage Status](https://coveralls.io/repos/JustAdam/streamingtwitter/badge.png)](https://coveralls.io/r/JustAdam/streamingtwitter) [![GoDoc](https://godoc.org/github.com/JustAdam/streamingtwitter?status.svg)](https://godoc.org/github.com/JustAdam/streamingtwitter)

Go package to provide access to Twitter's streaming and REST API.

This is a early version which is no doubt missing lots of functionality (and tests), so just open an issue and/or help out :)


Example clients and how to use the code can be found in the `cmd/` directory.


Quick start
-----------

 	$ go get github.com/JustAdam/streamingtwitter

	$ cd cmd/

	$ mv tokens.json.sample tokens.json

	$ vim tokens.json
		Add your Twitter API token and secret under "App".

	$ go run simple/simple.go
		Follow the instructions to grant access to your Twitter app
