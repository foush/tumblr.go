package tumblrapi

import (
	"testing"
	"net/url"
	"net/http"
)

type testEndpoint struct {
	fn func()
}

var testEndpoints []testEndpoint = []testEndpoint{}

func getDashString(posts...Post) string {
	return jsonStringify(map[string]interface{}{
		"response": map[string]interface{}{
			"posts": posts,
		},
	})
}

func TestPaginateDashboardByOffset(t *testing.T) {
	p3 := Post{PostRef: PostRef{MiniPost: MiniPost{Id: 12345}}}
	c := newTestClient(getDashString(Post{}, Post{}, p3), nil)
	c.confirmExpectedSet = expectClientCallParams(t, "GetDashboard", http.MethodGet, "/user/dashboard", url.Values{})
	response, err := GetDashboard(c, url.Values{})
	if err != nil {
		t.Error("Unexpected error getting dashboard", err)
	}
	if response == nil {
		t.Error("Dashboard response was nil")
	}
	if response.byOffset || response.bySince {
		t.Error("Initial dashboard request should not be marked as paginated by `offset` or `sinceId`")
	}
	if len(response.Posts) != 3 {
		t.Error("Expected a dashboard with 3 posts")
	}
	c.response.body = []byte(getDashString())
	c.confirmExpectedSet = expectClientCallParams(t, "GetDashboard", http.MethodGet, "/user/dashboard", url.Values{
		"offset": []string{"3"},
	})
	nextResponse, err := response.NextByOffset()
	if err != nil {
		t.Error("Unexpected error getting next page of dashboard")
	}
	if response.params.Get("offset") == nextResponse.params.Get("offset") {
		t.Error("Next page should not modify params of previous result")
	}
	if len(nextResponse.Posts) > 0 {
		t.Fatal("Next dashboard should have 0 posts")
	}
	_, err = nextResponse.NextByOffset()
	if err != NoNextPageError {
		t.Error("Empty response should prevent next page")
	}
	_, err = nextResponse.NextBySinceId()
	if err != MixedPaginationMethodsError {
		t.Error("Changing pagination should generate error")
	}
}

func TestPaginateDashboardById(t *testing.T) {
	p3 := Post{PostRef: PostRef{MiniPost: MiniPost{Id: 12345}}}
	c := newTestClient(getDashString(Post{}, Post{}, p3), nil)
	c.confirmExpectedSet = expectClientCallParams(t, "GetDashboard", http.MethodGet, "/user/dashboard", url.Values{})
	response, err := GetDashboard(c, url.Values{})
	if err != nil {
		t.Error("Unexpected error getting dashboard", err)
	}
	if response == nil {
		t.Error("Dashboard response was nil")
	}
	if response.byOffset || response.bySince {
		t.Error("Initial dashboard request should not be marked as paginated by `offset` or `sinceId`")
	}
	if len(response.Posts) != 3 {
		t.Error("Expected a dashboard with 3 posts")
	}
	c.response.body = []byte(getDashString())
	c.confirmExpectedSet = expectClientCallParams(t, "GetDashboard", http.MethodGet, "/user/dashboard", setParamsUint(p3.Id, url.Values{}, "since_id"))
	nextResponse, err := response.NextBySinceId()
	if err != nil {
		t.Error("Unexpected error getting next page of dashboard")
	}
	if response.params.Get("since_id") == nextResponse.params.Get("since_id") {
		t.Error("Next page should not modify params of previous result")
	}
	if len(nextResponse.Posts) > 0 {
		t.Fatal("Next dashboard should have 0 posts")
	}
	_, err = nextResponse.NextByOffset()
	if err != MixedPaginationMethodsError {
		t.Error("Empty response should prevent next page")
	}
	_, err = nextResponse.NextBySinceId()
	if err != NoNextPageError {
		t.Error("Changing pagination should generate error")
	}
}
