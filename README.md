# Tumblr API Go Client

[![Build Status](https://travis-ci.org/foush/tumblr.go.svg?branch=master)](https://travis-ci.org/foush/tumblr.go) [![GoDoc](https://godoc.org/github.com/foush/tumblr.go?status.svg)](https://godoc.org/github.com/foush/tumblr.go)

This is the Tumblr API Golang client

## Installation

Run `go get github.com/foush/tumblr.go`

## Usage

The mechanics of this library send HTTP requests through a `ClientInterface` object. There is intentionally no concrete client defined in this library to allow for maximum flexibility. There is [a separate repository](https://github.com/foush/tumblrclient.go) with a client implementation and convenience methods if you do not require a custom client behavior.

