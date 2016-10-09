package tumblrapi

import (
	"net/url"
	"encoding/json"
)

type Likes struct {
	client ClientInterface
	response *Response
	parsedPosts []PostInterface
	Posts []MiniPost `json:"liked_posts"`
	TotalLikes uint64 `json:"liked_count"`
}

// URL values can include:
// 	limit (int)
//	offset (int)
//	before (timestamp)
//	after (timestamp)
func GetLikes(client ClientInterface, params url.Values) (*Likes, error) {
	response, err := client.GetWithParams("/user/likes", params)
	if err != nil {
		return nil, err
	}

	result := struct {
		Response Likes `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	result.Response.client = client
	result.Response.response = &response
	return &result.Response, nil
}

func doLike(client ClientInterface, path string, postId uint64, reblogKey string) error {
	_, err := client.PostWithParams(path, setPostId(postId, url.Values{
		"reblog_key": []string{reblogKey},
	}))
	return err
}

func LikePost(client ClientInterface, postId uint64, reblogKey string) error {
	return doLike(client, "/user/like", postId, reblogKey)
}

func UnlikePost(client ClientInterface, postId uint64, reblogKey string) error {
	return doLike(client, "/user/unlike", postId, reblogKey)
}

func (l *Likes)Full() ([]PostInterface, error) {
	var err error = nil
	if l.parsedPosts == nil {
		r := struct {
			Response struct{
					 Posts []PostInterface `json:"liked_posts"`
				 } `json:"response"`
		}{}
		r.Response.Posts = makePostsFromMinis(l.Posts, l.client)
		//fmt.Println(string(p.response.body))
		if err = json.Unmarshal(l.response.body, &r); err != nil {
			l.parsedPosts = []PostInterface{}
		} else {
			l.parsedPosts = r.Response.Posts
		}
	}
	return l.parsedPosts, err
}
