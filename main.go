package main

import (
	"fmt"
	"os"

	"github.com/aybabtme/rgbterm"
	"github.com/codegangsta/cli"
	"github.com/ianferguson/coven/posts"
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
		posts, err := posts.Get(12)
		if err != nil {
			panic(err)
		}

		for _, post := range posts {
			fmt.Printf("%v:\n", post.Summary())
			fmt.Printf("\tarticle: %v\n", blue(post.URL))
			fmt.Printf("\tcomments(%v): %v\n\n", post.CommentCount, blue(post.Comments))
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
