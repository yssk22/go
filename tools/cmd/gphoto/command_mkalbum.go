// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/speedland/go/services/google/photo"
)

var commandMkalbum = cli.Command{
	Name:  "mkalbum",
	Usage: "create an album",
	Action: func(c *cli.Context) {
		mkalbum(c)
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "public",
			Usage: "set public album",
		},
	},
}

func mkalbum(c *cli.Context) {
	client := NewClient()
	args := c.Args()
	if len(args) == 0 {
		fmt.Println("No album name is specified.")
		return
	}
	for _, arg := range args {
		a, err := client.CreateAlbum("", newAlbum(c, arg))
		if err != nil {
			printError(err)
			continue
		}
		fmt.Printf("Album %s was created.\n", a.ID)
	}
}

func newAlbum(c *cli.Context, title string) *photo.Album {
	a := photo.NewAlbum()
	a.Title = title
	if c.IsSet("public") {
		a.Access = photo.AccessPublic
	}
	return a
}
