package tumblrapi

import (
	"testing"
	"net/http"
	"net/url"
	"reflect"
	"fmt"
	"errors"
)

func TestPostRefLike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "reblog-key"
	ref := PostRef{
		client: client,
		MiniPost: MiniPost{
			Id: postId,
			ReblogKey: reblogKey,
		},
	}
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"PostRef.Like",
		http.MethodPost,
		"/user/like",
		params,
	)
	ref.Like()
}

func TestPostRefUnlike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "reblog-key"
	ref := PostRef{
		client: client,
		MiniPost: MiniPost{
			Id: postId,
			ReblogKey: reblogKey,
		},
	}
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"PostRef.Like",
		http.MethodPost,
		"/user/unlike",
		params,
	)
	ref.Unlike()
}

func TestMakePostFromType(t *testing.T) {
	testCases := map[string]string{
		"quote": "QuotePost",
		"chat": "ChatPost",
		"photo": "PhotoPost",
		"text": "TextPost",
		"link": "LinkPost",
		"answer": "AnswerPost",
		"audio": "AudioPost",
		"video": "VideoPost",
	}
	classPrefix := "*tumblrapi."
	for postType,postClass := range testCases {
		post,err := makePostFromType(postType)
		if err != nil {
			t.Errorf("Unexpected error creating post of type `%s`", postType)
		}
		postClass = classPrefix + postClass
		if actual := reflect.TypeOf(post).String(); actual != postClass {
			t.Errorf("Expected `%s` type to generate struct type `%s`, got `%s` instead", postType, postClass, actual)
		}
	}
	// test default case
	_,err := makePostFromType("")
	if err == nil {
		t.Fatal("Unexpected type should generate an error")
	}
}

func TestStringifyPost (t *testing.T) {
	post := Post{}
	if post.String() != jsonStringify(post) {
		t.Fatal("Post stringify does not conform to expected JSON output")
	}
}

func TestPostDynamicAccessor (t *testing.T) {
	post := Post{}
	post.Id = 1986
	if _, err := post.GetProperty("DoesNotExistProperty"); err == nil {
		t.Fatal("Dynamic accessor should error on property that does not exist")
	}
	actual, err := post.GetProperty("Id")
	if err != nil {
		t.Fatal("Dynamic accessor incorrectly errored")
	}
	if fmt.Sprintf("%v", actual) != "1986" {
		t.Fatalf("Dynamic accessor does not return the proper value; expected %d got %v", post.Id, actual)
	}
}

func TestQueryPostsReturnsClientError (t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("", clientErr)
	if _, err := queryPosts(client, "", "", url.Values{}); err == nil {
		t.Fatal("Client error should be returned")
	}
}

func TestQueryPostsReturnsJsonError (t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := queryPosts(client, "", "", url.Values{}); err == nil {
		t.Fatal("JSON Unmarshal error should be returned")
	}
}

func TestQueryPostsSuccess (t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	path := "/blog/%s/something"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"queryPosts",
		http.MethodGet,
		blogPath(path, blogName),
		params,
	)
	response, err := queryPosts(
		client,
		path,
		blogName,
		params,
	)
	if err != nil {
		t.Fatal("Posts should have been returned")
	}
	if string(response.response.body) != string(client.response.body) {
		t.Fatal("Response should match client's response")
	}
}

func TestGetPosts (t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts", blogName),
		params,
	)
	if _, err := GetPosts(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetQueue (t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/queue", blogName),
		params,
	)
	if _, err := GetQueue(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetDrafts (t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/draft", blogName),
		params,
	)
	if _, err := GetDrafts(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetSubmissions (t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/submission", blogName),
		params,
	)
	if _, err := GetSubmissions(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}