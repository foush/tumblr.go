package tumblrapi

import (
	"net/url"
	"strconv"
	"encoding/json"
)

type FollowingList struct {
	client ClientInterface
	Total uint32 `json:"total_blogs"`
	Blogs []Blog `json:"blogs"`
	Offset uint
	Limit uint
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

// get list of blogs this user follows
func GetFollowing(client ClientInterface, offset, limit uint) (*FollowingList, error) {
	params := url.Values{
		"limit": []string{strconv.FormatUint(uint64(limit), 10)},
		"offset": []string{strconv.FormatUint(uint64(offset), 10)},
	}
	result, err := client.GetWithParams("/user/following", params)
	if err != nil {
		return nil, err
	}
	response := struct{
		Response FollowingList `json:"response"`
	}{
		Response: FollowingList{
			client: client,
			Limit: limit,
			Offset: offset,
		},
	}
	if err = json.Unmarshal(result.body, &response); err != nil {
		return nil, err
	}
	if offset < 1 {
		response.Response.Limit = uint(len(response.Response.Blogs))
	}
	return &response.Response, nil
}


func (f *FollowingList) Next() (*FollowingList, error) {
	offset := f.Offset + f.Limit
	if offset >= uint(f.Total) {
		return nil, NoNextPageError
	}
	return GetFollowing(f.client, offset, f.Limit)
}

func (f *FollowingList) Prev() (*FollowingList, error) {
	if f.Offset <= 0 {
		return nil, NoPrevPageError
	}
	newOffset := f.Offset - f.Limit
	if newOffset < 0 {
		newOffset = 0
	}
	return GetFollowing(f.client, newOffset, f.Offset)
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

// follow a blog
func Follow(client ClientInterface, blogName string) error {
	_, err := client.PostWithParams("/user/follow", url.Values{
		"url": []string{normalizeBlogName(blogName)},
	})
	return err
}

// unfollow a blog
func Unfollow(client ClientInterface, blogName string) error {
	_, err := client.PostWithParams("/user/unfollow", url.Values{
		"url": []string{normalizeBlogName(blogName)},
	})
	return err
}
