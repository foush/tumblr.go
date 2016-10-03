# Tumblr API Go Client

This is the Tumblr API Golang client

## Installation

Run `go get github.com/foush/tumblrapi`

This project utilizes an external OAuth1 library you can find at [github.com/dghubble/oauth1](github.com/dghubble/oauth1) 

## Usage

After installing you can create a tumblr api client with either your consumer credentials or with both consumer and a user's credentials

```
	client := tumblrapi.NewClient(
		"CONSUMER KEY",
		"CONSUMER SECRET",
	)
	// or
    client := tumblr_go.NewClientWithToken(
		"CONSUMER KEY",
		"CONSUMER SECRET",
		"USER TOKEN",
		"USER TOKEN SECRET",
    )
```

The client is necessary to pass in to any API method. If you'd rather utilize your own OAuth client, you can create a wrapper around it that implements the `ClientInterface` and pass that in to the API methods instead.

# Methods

## Posts

### CreatePost
* **params:** *client ClientInterface, name string, params url.Values
* **returns:** *uint64, error*
Creates a post if the parameters validate and returns the post id, else an error

// Edit a given post, returns nil if successful, error on failure
### EditPost
* **params:** *client ClientInterface, name string, postId uint64, params url.Values*
* **returns:** *error*
Edits a given post. Returns nil if successful, an error on failure

// Reblog a given post to the given blog, returns the reblog's post id if successful, else the error
### ReblogPost
* **params:** *client ClientInterface, name string, postId uint64, reblogKey string, params url.Values*
* **returns:** *uint64, error*
Creates a reblog of the post at `postId` if the parameters validate and returns the reblog post id, else an error

### DeletePost
* **params:** *client ClientInterface, name string, postId uint64
* **returns:** *error*
Deletes the post at `postId` from `{name}.tumblr.com`. Returns nil if successful, an error on failure

## Tumblelog

### GetBlog
* **params:** *client ClientInterface, name string
* **returns:** *\*Blog, error*
Retrieves information about the blog `{name}.tumblr.com` 

### GetAvatar
* **params:** *client ClientInterface, name string
* **returns:** *string, error*
Retrieves the URI for `{name}.tumblr.com`'s avatar

### GetFollowers
* **params:** *client ClientInterface, name string*
* **returns:** *\*FollowerList, error*
Retrieves a list of followers for `{name}.tumblr.com`.

### GetPosts
* **params:** *client ClientInterface, name string, params url.Values*
* **returns:** *\*Posts, error*
Queries the API for a list of posts from blog at `{name}.tumblr.com`. You may pass modification query params as described by the API documentation

## GetQueue
* **params:** *client ClientInterface, name string, params url.Values*
* **returns:** *\*Posts, error*
Queries the API for a list of queued posts from blog at `{name}.tumblr.com`. You may pass modification query params as described by the API documentation

## GetDrafts
* **params:** *client ClientInterface, name string, params url.Values*
* **returns:** *\*Posts, error*
Queries the API for a list of draft posts from blog at `{name}.tumblr.com`. You may pass modification query params as described by the API documentation

## GetSubmissions
* **params:** *client ClientInterface, name string, params url.Values*
* **returns:** *\*Posts, error*
Queries the API for a list of posts submitted to `{name}.tumblr.com`. You may pass modification query params as described by the API documentation

