package tumblrapi

import (
	"encoding/json"
)

type User struct {
	Following uint32 `json:"following"`
	DefaultPostFormat string `json:"default_post_format"`
	Name string `json:"name"`
	Likes uint64 `json:"likes"`
	Blogs []ShortBlog `json:"blogs"`
}

func GetUserInfo(client ClientInterface) (*User, error) {
	response, err := client.Get("/user/info")
	if err != nil {
		return nil, err
	}
	result := struct{
		Response struct {
			User User `json:"user"`
		} `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	return &result.Response.User, nil
}


//func Follow(client ClientInterface, blogName string) {
//
//}
//
//func Unfollow(client ClientInterface, blogName string) {
//
//}

//
//func (l *Likes)GetAll() {
//	if l.parsedPosts == nil {
//		r := struct {
//			Response struct {
//					 Posts []PostInterface `json:"posts"`
//				 } `json:"response"`
//		}{}
//		r.Response.Posts = makePostsFromMinis(l.Posts, l.client)
//		//fmt.Println(string(p.response.body))
//		if err := json.Unmarshal(l.response.body, &r); err != nil {
//			l.parsedPosts = []PostInterface{}
//		} else {
//			l.parsedPosts = r.Response.Posts
//		}
//	}
//	return l.parsedPosts
//}