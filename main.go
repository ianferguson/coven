package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ianferguson/coven/posts"
)

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
			fmt.Printf("\tarticle: %v\n", post.URL)
			fmt.Printf("\tcomments(%v): %v\n\n", post.CommentCount, post.Comments)
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%v", err)
	}
}
