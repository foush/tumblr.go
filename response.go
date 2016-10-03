package tumblrapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// API Response structure which we'll use to pass back to behaviors
type Response struct {
	body []byte
	Headers http.Header
	Meta map[string]interface{}
	Result map[string]interface{}
	Errors map[string]interface{}
}

// Create a response object from the body bytestream and the headers structure
func NewResponse(body []byte, headers http.Header) *Response {
	return &Response{body: body, Headers: headers}
}

// Get the raw response body
func (r *Response) GetBody() []byte {
	return r.body
}

// Utility function for populating the Response's fields
func (r *Response) PopulateFromBody() error {
	if len(r.body) < 1 {
		return errors.New("Unable to populate from empty body")
	}
	// already populated, don't do again
	if r.Meta != nil || r.Result != nil || r.Errors != nil {
		return nil
	}
	data := map[string]interface{}{}
	e := json.Unmarshal(r.body, &data)
	if e != nil {
		return e
	}
	if value,ok := data["meta"]; ok {
		r.Meta, ok = value.(map[string]interface{})
	}
	if value,ok := data["response"]; ok {
		r.Result, ok = value.(map[string]interface{})
	}
	if value,ok := data["errors"]; ok {
		r.Errors, ok = value.(map[string]interface{})
	}
	return nil
}