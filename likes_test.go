package tumblrapi

import (
	"testing"
	"errors"
	"net/url"
	"net/http"
)

func TestGetLikesClientErrorReturnsError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("", clientErr)
	if _, err := GetLikes(client, url.Values{}); err != clientErr {
		t.Fatal("Client error should be returned")
	}
}

func TestGetLikesJSONErrorReturnsError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetLikes(client, url.Values{}); err == nil {
		t.Fatal("JSON unmarshal error should be returned")
	}
}

func TestGetLikesSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetLikes",
		http.MethodGet,
		"/user/likes",
		params,
	)
	if response, err := GetLikes(client, params); err != nil || response == nil {
		t.Fatal("Request should succeed")
	} else if response.response == nil {
		t.Fatal("Response should be set")
	} else if response.client != client {
		t.Fatal("Response should set client")
	}
}

func TestDoLikeError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	path := ""
	var postId uint64 = 1986
	reblogKey := ""
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doLike",
		http.MethodPost,
		path
	)
	doLike(client, path, postId, reblogKey)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetLikes",
		http.MethodGet,
		"/user/likes",
		params,
	)
	if response, err := GetLikes(client, params); err != nil || response == nil {
		t.Fatal("Request should succeed")
	} else if response.response == nil {
		t.Fatal("Response should be set")
	} else if response.client != client {
		t.Fatal("Response should set client")
	}
}
