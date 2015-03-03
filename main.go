package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ianferguson/coven/posts"
	"github.com/ianferguson/rgbterm"
)

func blue(s string) string {
	return rgbterm.FgString(s, 6, 69, 173)
}

func main() {
	app := cli.NewApp()
	app.Name = "coven"
	app.Usage = "download and review stories from coven"
	app.Email = "ian@labmarie.com"
	app.Version = "0.1"
	app.Action = func(c *cli.Context) {
		p, err := posts.Newest(12)
		if err != nil {
			panic(err)
		}
		printPosts("Newest Posts", p)

		p, err = posts.MostDiscussed(4)
		if err != nil {
			panic(err)
		}
		printPosts("Most Discussed", p)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
	}
}

func printPosts(title string, p posts.Posts) {
	fmt.Printf("%v\n-------------------------\n", title)
	for _, post := range p {
		fmt.Printf("%v:\n", post.Summary())
		fmt.Printf("\tarticle: %v\n", blue(post.URL))
		fmt.Printf("\tcomments(%v): %v\n\n", post.CommentCount, blue(post.Comments))
	}
}
