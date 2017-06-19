// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"fmt"
	"time"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/speedland/go/services/google/photo"
)

var commandLs = cli.Command{
	Name:  "ls",
	Usage: "list albums or photos in albums",
	Action: func(c *cli.Context) {
		ls(c)
	},
}

func ls(c *cli.Context) {
	client := NewClient()
	args := c.Args()
	if len(args) == 0 {
		printAlbums(client)
	} else {
		for _, arg := range args {
			if strings.Index(arg, "/") > 0 {
				tmp := strings.Split(arg, "/")
				printMediumByID(client, tmp[0], tmp[1])
			} else {
				printMediaByAlbumID(client, arg)
			}
		}
	}
}

func printAlbums(client *photo.Client) {
	albums, err := client.ListAlbums("")
	if err != nil {
		panic(err)
	}
	table := make([][]string, 0)
	table = append(table, []string{
		"ID",
		"UpdatedAt",
		"Access",
		"NumPhotos",
		"BytesUsed",
		"Title",
	})
	for _, a := range albums {
		table = append(table, []string{
			a.ID,
			a.UpdatedAt.Local().Format(time.RFC3339),
			string(a.Access),
			fmt.Sprintf("%d/%d", a.NumPhotos, (a.NumPhotos + a.NumPhotosRemaining)),
			fmt.Sprintf("%d", a.BytesUsed),
			a.Title,
		})
	}
	printTable(table)
}

func printMediaByAlbumID(client *photo.Client, albumID string) {
	_media, err := client.ListMedia("", albumID)
	if err != nil {
		panic(err)
	}
	table := make([][]string, 0)
	table = append(table, []string{
		"ID",
		"UpdatedAt",
		"Size",
		"Status",
		"Title",
	})
	for _, m := range _media {
		table = append(table, []string{
			fmt.Sprintf("%s/%s", m.AlbumID, m.ID),
			m.UpdatedAt.Local().Format(time.RFC3339),
			fmt.Sprintf("%d", m.Size),
			m.VideoStatus,
			m.Title,
		})
	}
	printTable(table)
}

func printMediumByID(client *photo.Client, albumID string, mediaID string) {
	media, err := client.GetMedia("", albumID, mediaID)
	if err != nil {
		panic(err)
	}
	table := make([][]string, 0)
	table = append(table, []string{
		"URL",
		"Type",
		"Width",
		"Height",
		"Status",
	})
	for _, m := range media.Contents {
		table = append(table, []string{
			m.URL,
			m.Type,
			fmt.Sprintf("%d", m.Width),
			fmt.Sprintf("%d", m.Height),
		})
	}
	printTable(table)
}
