package tumblrapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"net/url"
	"strconv"
)

// Representation of a list of Posts
type Posts struct {
	response Response
	parsedPosts []PostInterface
	Posts []MiniPost `json:"posts"`
	TotalPosts int64 `json:"total_posts"`
}

// Method to retrieve fully fleshed post data from stubs and cache result
func (p *Posts) All() ([]PostInterface, error) {
	var err error = nil
	if p.parsedPosts == nil {
		posts := []PostInterface{}
		for _, mini := range p.Posts {
			post, _ := makePostFromType(mini.Type)
			posts = append(posts, post)
		}
		r := struct {
			Response struct {
					 Posts []PostInterface `json:"posts"`
				 } `json:"response"`
		}{}
		r.Response.Posts = posts
		//fmt.Println(string(p.response.body))
		if err = json.Unmarshal(p.response.body, &r); err != nil {
			p.parsedPosts = []PostInterface{}
		} else {
			p.parsedPosts = posts
		}
	}
	return p.parsedPosts, err
}

// Method to retrieve a single Post entity at a given index; returns nil if index is out of bounds
func (p *Posts) Get(index int) (PostInterface) {
	if index < 0 {
		return nil
	}
	if posts,err := p.All(); err == nil {
		if index >= len(posts) {
			return nil
		}
		return posts[index]
	}
	return nil
}

// The basics for what is needed in a Post
type MiniPost struct {
	Id uint64 `json:"id"`
	Type string `json:"type"`
}

// The common fields on any post, no matter what type
type Post struct {
	PostInterface
	MiniPost
	BlogName string `json:"blog_name"`
	Body string `json:"body"`
	CanLike bool `json:"can_like"`
	CanReblog bool `json:"can_reblog"`
	CanReply bool `json:"can_reply"`
	CanSendInMessage bool `json:"can_send_in_message"`
	Caption string `json:"caption"`
	Date string `json:"date"`
	DisplayAvatar bool `json:"display_avatar"`
	Followed bool `json:"followed"`
	Format string `json:"format"`
	Highlighted []interface{} `json:"highlighted"`
	Liked bool `json:"liked"`
	NoteCount uint64 `json:"note_count"`
	PermalinkUrl string `json:"permalink_url"`
	PostUrl string `json:"post_url"`
	Reblog struct {
		 Comment string `json:"comment"`
		 TreeHTML string `json:"tree_html"`
	       } `json:"reblog"`
	ReblogKey string `json:"reblog_key"`
	RecommendedColor string `json:"recommended_color"`
	RecommendedSource bool `json:"recommended_source"`
	ShortUrl string `json:"short_url"`
	Slug string `json:"slug"`
	SourceTitle string `json:"source_title"`
	SourceUrl string `json:"source_url"`
	State string `json:"state"`
	Summary string `json:"summary"`
	Tags []string `json:"tags"`
	Timestamp uint64 `json:"timestamp"`
	TrackName string `json:"track_name,omitempty"`
	Trail []ReblogTrailItem `json:"trail"`
}

// Post substructure
type ReblogTrailItem struct {
	Blog Blog `json:"blog"`
	Content string `json:"content"`
	ContentRaw string `json:"content_raw"`
	IsCurrentItem bool `json:"is_current_item"`
	Post struct {
		     // sometimes an actual int, sometimes a numeric string, always a headache
		     Id interface{} `json:"id"`
	     } `json:"post"`
}

// PostInterface for use in typed structures which could contain any of the below subtypes
type PostInterface interface {
	GetProperty(key string) (interface{}, error)
	GetSelf() (*Post)
}

// Post subtype
type QuotePost struct {
	Post
	Source string `json:"source,omitempty"`
	Text string `json:"text"`
}

// Post subtype
type ChatPost struct {
	Post
	Dialog []struct{
		Label string `json:"label"`
		Name string `json:"name"`
		Phrase string `json:"phrase"`
	} `json:"dialog"`
}

// Post subtype
type TextPost struct {
	Post
	Title string `json:"title"`
}

// Post subtype
type LinkPost struct {
	Post
	Description string `json:"description"`
	Excerpt string `json:"excerpt"`
	LinkAuthor string `json:"link_author"`
	Title string `json:"title"`
	Url string `json:"url"`
}

// Post subtype
type AnswerPost struct {
	Post
	Answer string `json:"answer"`
	AskingName string `json:"asking_name"`
	AskingUrl string `json:"asking_url"`
	Publisher string `json:"publisher"`
	Question string `json:"question"`
}

// Post subtype
type AudioPost struct {
	Post
	AlbumArt string `json:"album_art"`
	Artist string `json:"artist"`
	AudioSourceUrl string `json:"audio_source_url"`
	AudioType string `json:"audio_type"`
	AudioUrl string `json:"audio_url"`
	Embed string `json:"embed"`
	Player string `json:"player"`
	Plays uint64 `json:"plays"`
}

// Post subtype
type VideoPost struct {
	Post
	Html5Capable bool `json:"html5_capable"`
	PermalinkUrl string `json:"permalink_url"`
	Players []struct {
		EmbedCode string `json:"embed_code"`
		Width interface{} `json:"width"`
	} `json:"player"`
	ThumbnailHeight uint32 `json:"thumbnail_height"`
	ThumbnailUrl string `json:"thumbnail_url"`
	ThumbnailWidth uint32 `json:"thumbnail_width"`
	Video map[string]struct {
		Height uint32 `json:"height"`
		Width uint32 `json:"width"`
		VideoId string `json:"video_id"`
	} `json:"video"`
	VideoType string `json:"video_type"`
}

// Post subtype
type PhotoPost struct {
	Post
	ImagePermalink string `json:"image_permalink"`
	Photos []Photo `json:"photos"`
}

// Photo post substructure
type Photo struct {
	AltSizes []PhotoSize `json:"alt_sizes"`
	Caption string `json:"caption"`
	OriginalSize PhotoSize `json:"original_size"`
}

// Photo substructure
type PhotoSize struct {
	Height uint32 `json:"height"`
	Width uint32 `json:"width"`
	Url string `json:"url"`
}

// Convenience method for ease of use
func (p *Post) String() string {
	return jsonStringify(*p)
}

// Convenience method for easy retrieval of one-off values
func (p *Post) GetProperty(key string) (interface{},error) {
	if field,exists := reflect.TypeOf(p).Elem().FieldByName(key); exists {
		return reflect.ValueOf(p).Elem().FieldByIndex(field.Index), nil
	}
	return nil, errors.New(fmt.Sprintf("Property %s does not exist", key))
}

// Useful for converting a PostInterface into a Post
func (p *Post) GetSelf() (*Post) {
	return p
}

// helper method for querying a given path which should return a list of posts
func queryPosts(client ClientInterface, path, name string, params url.Values) (*Posts, error) {
	response, err := client.GetWithParams(blogPath(path, name), params)
	if err != nil {
		return nil, err
	}
	posts := struct {
		Response Posts `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &posts); err == nil {
		posts.Response.response = response
		// store
		return &posts.Response, nil
	}
	return nil, err
}

// Retrieve a blog's posts, in the API docs you can find how to filter by ID, type, etc
func GetPosts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts", name, params)
}

// Retrieve a blog's Queue
func GetQueue(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/queue", name, params)
}

// Retrieve a blog's drafts
func GetDrafts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/draft", name, params)
}

// Retrieve a blog's submsisions
func GetSubmissions(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/submission", name, params)
}

// Util method for decoding the response
func doPost(client ClientInterface, path, name string, params url.Values) (uint64, error) {
	response, err := client.PostWithParams(blogPath(path, name), params)
	if err != nil {
		return 0, err
	}
	post := struct {
		Response struct{
				 Id uint64 `json:"id"`
			 } `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &post); err == nil {
		return post.Response.Id, nil
	}
	return 0, err
}

// Create a post, return the ID on success, error on failure
func CreatePost(client ClientInterface, name string, params url.Values) (uint64, error) {
	return doPost(client, "/blog/%s/post", name, params)
}

// Edit a given post, returns nil if successful, error on failure
func EditPost(client ClientInterface, name string, postId uint64, params url.Values) error {
	params.Set("id", strconv.FormatUint(postId, 10))
	_, err := client.PostWithParams(blogPath("/blog/%s/post/edit", name), params)
	return err
}

// Reblog a given post to the given blog, returns the reblog's post id if successful, else the error
func ReblogPost(client ClientInterface, name string, postId uint64, reblogKey string, params url.Values) (uint64, error) {
	params.Set("id", strconv.FormatUint(postId, 10))
	params.Set("reblog_key", reblogKey)
	return doPost(client, "/blog/%s/post/reblog", name, params)
}

// Delete a given blog's post by ID, nil if successful, error on failure
func DeletePost(client ClientInterface, name string, postId uint64) error {
	_, err := client.PostWithParams(blogPath("/blog/%s/post/delete", name), url.Values{"id": []string{strconv.FormatUint(postId, 10)}})
	return err
}

// Helper function
func GetText(ptr PostInterface) (*TextPost) {
	if post,ok := ptr.(*TextPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetAudio(ptr PostInterface) (*AudioPost) {
	if post,ok := ptr.(*AudioPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetVideo(ptr PostInterface) (*VideoPost) {
	if post,ok := ptr.(*VideoPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetChat(ptr PostInterface) (*ChatPost) {
	if post,ok := ptr.(*ChatPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetAnswer(ptr PostInterface) (*AnswerPost) {
	if post,ok := ptr.(*AnswerPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetLink(ptr PostInterface) (*LinkPost) {
	if post,ok := ptr.(*LinkPost); ok {
		return post
	}
	return nil
}

// Helper function
func GetQuote(ptr PostInterface) (*QuotePost) {
	if post,ok := ptr.(*QuotePost); ok {
		return post
	}
	return nil
}

// Utility function to create the proper instance of Post and return a reference to the generic interface
func makePostFromType(t string) (PostInterface, error) {
	switch t {
	case "quote":
		return &QuotePost{}, nil
	case "chat":
		return &ChatPost{}, nil
	case "photo":
		return &PhotoPost{}, nil
	case "text":
		return &TextPost{}, nil
	case "link":
		return &LinkPost{}, nil
	case "answer":
		return &AnswerPost{}, nil
	case "audio":
		return &AudioPost{}, nil
	case "video":
		return &VideoPost{}, nil
	}
	return &Post{}, errors.New(fmt.Sprintf("Unknown type %s", t))
}
