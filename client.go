package tumblrapi

import (
	"github.com/dghubble/oauth1"
	"net/http"
	"golang.org/x/net/context"
	"net/url"
	"io/ioutil"
	"errors"
	"strings"
	"encoding/json"
	"fmt"
)
const apiBase = "https://api.tumblr.com/v2/"

// The Tumblr API Client object
type Client struct {
	consumer *oauth1.Config
	user *oauth1.Token
	client *http.Client
}

// If you wish to use your own client, simply make sure it implements this interface
type ClientInterface interface {
	// Issue GET request to Tumblr API
	Get(endpoint string) (Response, error)
	// Issue GET request to Tumblr API with param values
	GetWithParams(endpoint string, params url.Values) (Response, error)
	// Issue POST request to Tumblr API
	Post(endpoint string) (Response, error)
	// Issue POST request to Tumblr API with param values
	PostWithParams(endpoint string, params url.Values) (Response, error)
	// Issue PUT request to Tumblr API
	Put(endpoint string) (Response, error)
	// Issue PUT request to Tumblr API with param values
	PutWithParams(endpoint string, params url.Values) (Response, error)
	// Issue DELETE request to Tumblr API
	Delete(endpoint string) (Response, error)
	// Issue DELETE request to Tumblr API with param values
	DeleteWithParams(endpoint string, params url.Values) (Response, error)
}

// Constructor with only the consumer key and secret
func NewClient(consumerKey string, consumerSecret string) *Client {
	c := Client{}
	c.SetConsumer(consumerKey, consumerSecret)
	return &c
}

// Constructor with consumer key/secret and user token/secret
func NewClientWithToken(consumerKey string, consumerSecret string, token string, tokenSecret string) *Client {
	c := NewClient(consumerKey, consumerSecret)
	c.SetToken(token, tokenSecret)
	return c
}

// Set consumer credentials, invalidates any previously cached client
func (c *Client) SetConsumer(consumerKey string, consumerSecret string) {
	c.consumer = oauth1.NewConfig(consumerKey, consumerSecret)
	c.client = nil
}

// Set user credentials, invalidates any previously cached client
func (c *Client) SetToken(token string, tokenSecret string) {
	c.user = oauth1.NewToken(token, tokenSecret)
	c.client = nil
}

// Issue GET request to Tumblr API
func (c *Client) Get(endpoint string) (Response, error) {
	return c.GetWithParams(endpoint, url.Values{})
}

// Issue GET request to Tumblr API with param values
func (c *Client) GetWithParams(endpoint string, params url.Values) (Response, error) {
	return getResponse(c.GetHttpClient().Get(createRequestURI(appendPath(apiBase,endpoint),params)))
}

// Issue POST request to Tumblr API
func (c *Client) Post(endpoint string) (Response, error) {
	return c.PostWithParams(endpoint, url.Values{});
}

// Issue POST request to Tumblr API with param values
func (c *Client) PostWithParams(endpoint string, params url.Values) (Response, error) {
	return getResponse(c.GetHttpClient().PostForm(appendPath(apiBase, endpoint), params))
}

// Issue PUT request to Tumblr API
func (c *Client) Put(endpoint string) (Response, error) {
	return c.PutWithParams(endpoint, url.Values{});
}

// Issue PUT request to Tumblr API with param values
func (c *Client) PutWithParams(endpoint string, params url.Values) (Response, error) {
	req, err := http.NewRequest("PUT", createRequestURI(appendPath(apiBase, endpoint), params), strings.NewReader(""))
	if err == nil {
		return getResponse(c.GetHttpClient().Do(req))
	}
	return Response{}, err
}

// Issue DELETE request to Tumblr API
func (c *Client) Delete(endpoint string) (Response, error) {
	return c.DeleteWithParams(endpoint, url.Values{});
}

// Issue DELETE request to Tumblr API with param values
func (c *Client) DeleteWithParams(endpoint string, params url.Values) (Response, error) {
	req, err := http.NewRequest("DELETE", createRequestURI(appendPath(apiBase, endpoint), params), strings.NewReader(""))
	if err == nil {
		return getResponse(c.GetHttpClient().Do(req))
	}
	return Response{}, err
}

// Retrieve the underlying HTTP client
func (c *Client) GetHttpClient() *http.Client {
	if c.consumer == nil {
		panic("Consumer credentials are not set")
	}
	if c.user == nil {
		c.SetToken("", "")
	}
	if c.client == nil {
		c.client = c.consumer.Client(context.TODO(), c.user)
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return c.client
}

// Helper function to ease appending path to a base URI
func appendPath(base string, path string) string {
	// if path starts with `/` shave it off
	if path[0] == '/' {
		path = path[1:]
	}
	return base + path
}

// Helper function to create a URI with query params
func createRequestURI(base string, params url.Values) string {
	if len(params) != 0 {
		base += "?" + params.Encode()
	}
	return base
}

// Standard way of receiving data from the API response
func getResponse(resp *http.Response, e error) (Response, error) {
	response := Response{}
	if e != nil {
		return response, e
	}
	defer resp.Body.Close()
	response.Headers = resp.Header
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return response, e
	}
	response = *NewResponse(body, resp.Header)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return response, errors.New(resp.Status)
	}
	return response, nil
}

// Helper function to JSON stringify a given value
func jsonStringify(b interface{}) string {
	out, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return string(out)
}
