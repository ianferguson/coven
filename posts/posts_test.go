package posts

import (
	"testing"
	"time"
)

func TestGetPosts(t *testing.T) {
	posts, _ := Get(10)

	if len(posts) != 10 {
		t.Errorf("found %v posts", len(posts))
	}
}

func TestSummary(t *testing.T) {
	now := time.Now()
	post := Post{
		Source:            "My New Blog",
		Title:             "Exciting Blog Post",
		CommentCount:      42,
		ExternalCreatedAt: &now,
	}

	summary := post.Summary()
	expected := "My New Blog: Exciting Blog Post, posted just a moment ago, 42 comments"
	if summary != expected {
		t.Errorf("post.Summary(%q) => %q, want %q", post, summary, expected)
	}
}

func TestPrettyPrint(t *testing.T) {
	// input is how many hours ago is being tested
	testCases := []struct {
		in  string
		out string
	}{
		{"-60h", "2 days ago"},
		{"-47h", "47 hours ago"},
		{"-90m", "90 minutes ago"},
		{"-30s", "just a moment ago"},
	}

	for _, testCase := range testCases {
		duration, _ := time.ParseDuration(testCase.in)
		time := time.Now().Add(duration)
		s := prettyPrint(&time)
		if s != testCase.out {
			t.Errorf("prettyPrint(%q) => %q, want %q", time, s, testCase.out)
		}
	}

}
