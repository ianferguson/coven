// Package posts contains methods for retrieving selections of posts from the coven api
package posts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	url = "http://api.coven.link/api/v1/posts"
)

// A Post contains information pertianing to a specific post made to coven
type Post struct {
	ID                int
	Position          int
	URL               string
	Title             string
	Comments          string
	CommentCount      int        `json:"comment_count"`
	ExternalCreatedAt *time.Time `json:"external_created_at,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	Source            string
}

// Summary returns a one line string summary of information relevant to the post
func (p *Post) Summary() string {
	posted := prettyPrint(p.ExternalCreatedAt)
	return fmt.Sprintf("%v, posted on %v %v", p.Title, p.Source, posted)
}

func prettyPrint(t *time.Time) string {
	since := time.Since(*t)
	switch {
	// this has some potential DST/timezone issues, but its just for display, so not a big issue
	case since.Hours() >= 48:
		count := int(since.Hours()) / 24
		return fmt.Sprintf("%v days ago", count)
	case since.Hours() > 2:
		return fmt.Sprintf("%v hours ago", int(since.Hours()))
	case since.Minutes() > 2:
		return fmt.Sprintf("%v minutes ago", int(since.Minutes()))
	}
	return "just a moment ago"
}

// GetAll retrieves all posts currently available via the coven api
func GetAll() ([]Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Get returns the most recent number of posts up to the number specified in the limit parameter
func Get(limit int) ([]Post, error) {
	posts, err := GetAll()
	if err != nil {
		return nil, err
	}

	if len(posts) < limit {
		return posts, nil
	}

	return posts[:limit], nil
}
