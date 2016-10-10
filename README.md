# Tumblr API Go Client

This is the Tumblr API Golang client

## Installation

Run `go get github.com/foush/tumblr.go`

This project utilizes an external OAuth1 library you can find at [github.com/dghubble/oauth1](github.com/dghubble/oauth1) 

## Usage

[![GoDoc](https://godoc.org/github.com/foush/tumblr.go?status.svg)](https://godoc.org/github.com/foush/tumblr.go)

The mechanics of this library send HTTP requests through a `ClientInterface` object. There is intentionally no concrete client defined in this library to allow for maximum flexibility. There is [a separate repository](https://github.com/foush/tumblrclient.go) with a client implementation and convenience methods if you do not require a custom client behavior.

