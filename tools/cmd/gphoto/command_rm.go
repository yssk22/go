// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/yssk22/go/services/google/photo"
)

var commandRm = cli.Command{
	Name:  "rm",
	Usage: "remove albums or photos in albums",
	Action: func(c *cli.Context) {
		rm(c)
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "force to remove even if any photos exist.",
		},
	},
}

func rm(c *cli.Context) {
	client := NewClient()
	args := c.Args()
	if len(args) == 0 {
		fmt.Println("No id is specified.")
		return
	}
	for _, arg := range args {
		if strings.Index(arg, "/") > 0 {
			tmp := strings.Split(arg, "/")
			_rmMedia(c, client, tmp[0], tmp[1])
		} else {
			_rmAlbum(c, client, arg)
		}
	}
}

func _rmAlbum(c *cli.Context, client *photo.Client, albumID string) {
	if !c.IsSet("force") {
		media, err := client.ListMedia("", albumID)
		if err != nil {
			printError(err)
			return
		}
		if media != nil && len(media) > 0 {
			printError(fmt.Errorf("The album still has photos. Use -f to force to remove."))
			return
		}
	}
	err := client.DeleteAlbum("", albumID)
	if err != nil {
		printError(err)
		return
	}
	fmt.Printf("Album %s was deleted.\n", albumID)
}

func _rmMedia(c *cli.Context, client *photo.Client, albumID string, mediaID string) {
	err := client.DeleteMedia("", albumID, mediaID)
	if err != nil {
		printError(err)
		return
	}
	fmt.Printf("Media %s was deleted.\n", mediaID)
}
