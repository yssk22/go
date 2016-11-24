package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/urfave/cli"
)

var classify = cli.Command{
	Name:  "classify",
	Usage: "classify a image file or URL",
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "id",
			Usage: "classifier id",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) != 1 {
			return fmt.Errorf("must specify a url or a file")
		}
		if strings.HasPrefix(args[0], "http://") || strings.HasPrefix(args[0], "https://") {
			params := &vr.ClassifyParams{
				ClassifierIDs: c.StringSlice("id"),
			}
			resp, err := client.ClassifyURL(context.Background(), args[0], params)
			if err != nil {
				return err
			}
			if len(resp.Images) == 0 {
				return fmt.Errorf("unexpected response: watson returned no images without errors")
			}
			resp.Images[0].Image = args[0]
			printImages(resp)
		}
		return nil
	},
}

func printImages(resp *vr.ClassifyResponse) {
	for _, image := range resp.Images {
		fmt.Printf("%s:\n", image.Image)
		for _, classifier := range image.Classifiers {
			fmt.Printf("\t%s:\n", classifier.Name)
			for _, class := range classifier.Classes {
				fmt.Printf("\t\t- %s: %f\n", class.Class, class.Score)
			}
		}
	}
}
