// Package posts contains methods for retrieving selections of posts from the coven api
package posts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
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

// Posts is a slice of Post pointers and associated methods for manipulating those Post structs
type Posts []*Post

// Limit returns the first n elements in a slice, or the entire slice, if len(slice) < n
func (s Posts) Limit(limit int) Posts {
	if len(s) < limit {
		return s
	}
	return s[:limit]
}

// post sorter takes a slice of Post's and a function to be used for sorting those posts
type postSorter struct {
	Posts
	by func(p1, p2 *Post) bool
}

// Len supports the sort.Interface methods needed to sort Posts
func (s *postSorter) Len() int {
	return len(s.Posts)
}

// Swap supports the sort.Interface methods needed to sort Posts
func (s *postSorter) Swap(i, j int) {
	s.Posts[i], s.Posts[j] = s.Posts[j], s.Posts[i]
}

// Less satisfies the sort.Interface via the provided function in postSorter
func (s *postSorter) Less(i, j int) bool {
	return s.by(s.Posts[i], s.Posts[j])
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
func GetAll() (Posts, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts Posts
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Newest returns the most recent number of posts up to the number specified in the limit parameter
// sorted by their date of posting descending
func Newest(limit int) (Posts, error) {
	posts, err := GetAll()
	if err != nil {
		return nil, err
	}

	sorter := func(p1, p2 *Post) bool {
		t1 := *p1.ExternalCreatedAt
		t2 := *p2.ExternalCreatedAt
		return t1.Sub(t2).Seconds() > 0
	}

	return sortPosts(posts, sorter).Limit(limit), nil
}

// MostDiscussed returns the most commented posts returned from coven.link including up to the
// number of posts specified by the limit parameter
func MostDiscussed(limit int) (Posts, error) {
	posts, err := GetAll()
	if err != nil {
		return nil, err
	}

	sorter := func(p1, p2 *Post) bool {
		return p1.CommentCount > p2.CommentCount
	}
	return sortPosts(posts, sorter).Limit(limit), nil
}

func sortPosts(p Posts, sortFunc func(p1, p2 *Post) bool) Posts {
	sorter := &postSorter{
		Posts: p,
		by:    sortFunc,
	}
	sort.Sort(sorter)
	return sorter.Posts
}
