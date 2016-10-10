package tumblrapi

import (
	"testing"
	"errors"
	"net/http"
	"net/url"
)

func TestBlog_String(t *testing.T) {
	blog := Blog{}
	if blog.String() != jsonStringify(blog) {
		t.Fatal("Blog JSON representation is incorrect")
	}
}

func TestGetBlogInfoClientErro(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	if _, err := GetBlogInfo(client, "name"); err != clientErr {
		t.Fatal("Blog info must return client error")
	}
}
func TestGetBlogInfoJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetBlogInfo(client, "name"); err == nil {
		t.Fatal("Blog info must return json unmarshal error")
	}
}
func TestGetBlogInfo(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"BlogInfo",
		http.MethodGet,
		blogPath("/blog/%s/info", blog),
		url.Values{},
	)
	if result, err := GetBlogInfo(client, blog); err != nil {
		t.Fatal("Blog info failed")
	} else if result.client != client {
		t.Fatal("Result must have client set.")
	}
}

func TestBlogRef_CreatePost(t *testing.T) {

}

func TestBlogRef_Follow(t *testing.T) {

}

func TestBlogRef_GetAvatar(t *testing.T) {

}

func TestBlogRef_GetDrafts(t *testing.T) {

}

func TestBlogRef_GetFollowers(t *testing.T) {

}

func TestBlogRef_GetInfo(t *testing.T) {

}

func TestBlogRef_GetPosts(t *testing.T) {

}

func TestBlogRef_GetQueue(t *testing.T) {

}

func TestBlogRef_ReblogPost(t *testing.T) {

}

func TestBlogRef_Unfollow(t *testing.T) {

}

func TestGetAvatar(t *testing.T) {

}

func TestNewBlogRef(t *testing.T) {

}
