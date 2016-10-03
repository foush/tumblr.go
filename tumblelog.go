package tumblrapi

import (
	"strings"
	"errors"
	"encoding/json"
	"fmt"
)

// Tumblelog struct
type Blog struct {
	Url string `json:"url"`
	Title string `json:"title"`
	Name string `json:"name"`
	Posts int64 `json:"posts"`
	Ask bool `json:"ask"`
	AskAnon bool `json:"ask_anon"`
	AskAnonPageTitle string `json:"ask_page_title"`
	CanSendFanMail bool `json:"can_send_fan_mail"`
	CanSubmit bool `json:"can_submit"`
	CanSubscribe bool `json:"can_subscribe"`
	Description string `json:"description"`
	Followed bool `json:"followed"`
	IsBlockedFromPrimary bool `json:"is_blocked_from_primary"`
	IsNSFW bool `json:"is_nsfw"`
	ShareLikes bool `json:"share_likes"`
	SubmissionPageTitle string `json:"submission_page_title"`
	Subscribed bool `json:"subscribed"`
	TotalPosts int64 `json:"total_posts"`
	Updated int64 `json:"updated"`
	Theme BlogTheme `json:"theme"`
}

// Tumblelog substructure
type BlogTheme struct {
	AvatarShape string `json:"avatar_shape"`
	BackgroundColor string `json:"background_color"`
	BodyFont string `json:"body_font"`
	// somtimes a single int, sometimes a space-separated series of int's (string)
	HeaderBounds interface{} `json:"header_bounds"`
	HeaderFocusHeight uint32 `json:"header_focus_height"`
	HeaderFocusWidth uint32 `json:"header_focus_width"`
	HeaderFullHeight uint32 `json:"header_full_height"`
	HeaderFullWidth uint32 `json:"header_full_width"`
	HeaderImage string `json:"header_image"`
	HeaderImageFocused string `json:"header_image_focused"`
	HeaderImageScaled string `json:"header_image_scaled"`
	HeaderStretch bool `json:"header_stretch"`
	LinkColor string `json:"link_color"`
	ShowAvatar bool `json:"show_avatar"`
	ShowDescription bool `json:"show_description"`
	ShowHeaderImage bool `json:"show_header_image"`
	ShowTitle bool `json:"show_title"`
	TitleColor string `json:"title_color"`
	TitleFont string `json:"title_font"`
	TitleFontWeight string `json:"title_font_weight"`
}

// Object from the lsit of followers response
type FollowerList struct {
	TotalUsers uint32 `json:"total_users"`
	Followers []Follower `json:"users"`
}

// FollowerList substructure
type Follower struct {
	Following bool `json:"following"`
	Name string `json:"name"`
	Updated int64 `json:"updated"`
	Url string `json:"url"`
}

// Convenience method
func (b *Blog) String() string {
	return jsonStringify(*b)
}

// Get information about a blog
func GetBlog(client ClientInterface, name string) (*Blog, error) {
	response, err := client.Get(blogPath("/blog/%s/info", name))
	if err != nil {
		return nil, err
	}
	blog := struct {
		Response struct {
			Blog Blog `json:"blog"`
		} `json:"response"`
	}{}
	//blog := blogResponse{}
	err = json.Unmarshal(response.body, &blog)
	if err != nil {
		return nil, err
	}
	return &blog.Response.Blog, nil
}

// Get Blog's Avatar URI
func GetAvatar(client ClientInterface, name string) (string, error) {
	response, err := client.Get(blogPath("/blog/%s/avatar", name))
	if err != nil {
		return "", err
	}
	if location := response.Headers.Get("Location"); len(location) > 0 {
		return location, nil
	}
	if err = response.PopulateFromBody(); err != nil {
		return "", err
	}
	if l, ok := response.Result["location"]; ok {
		if location, ok := l.(string); ok {
			return location, nil
		}
	}
	return "", errors.New("Unable to detect avatar location")
}

// Retrieve User's followers
func GetFollowers(client ClientInterface, name string) (*FollowerList, error) {
	response, err := client.Get(blogPath("/blog/%s/followers", name))
	if err != nil {
		return nil, err
	}
	followers := struct {
		Followers FollowerList `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &followers); err == nil {
		return &followers.Followers, nil
	}
	return nil, err
}

// Helper function to allow for less verbose code
func normalizeBlogName(name string) string {
	if !strings.Contains(name, ".") {
		name = fmt.Sprintf("%s.tumblr.com", name)
	}
	return name
}

// Expects path to contain a single %s placeholder to be substituted with the result of normalizeBlogName
func blogPath(path, name string) string {
	return fmt.Sprintf(path, normalizeBlogName(name))
}