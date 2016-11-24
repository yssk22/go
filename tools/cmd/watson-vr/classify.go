package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"

	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/speedland/go/x/xarchive/xzip"
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
		cli.Float64Flag{
			Name:  "threshold",
			Usage: "classifier threshold",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) == 0 {
			return fmt.Errorf("must specify a url or a file")
		}
		tempdir, err := ioutil.TempDir("", "watson-vr")
		if err != nil {
			return fmt.Errorf("error creating a temp directory: %v", err)
		}
		defer os.RemoveAll(tempdir)
		rawSources, sources := prepareSources(tempdir, args...)
		if len(sources) == 0 {
			return fmt.Errorf("no files can be classified")
		}
		defer func() {
			for _, v := range rawSources {
				v.Close()
			}
		}()
		params := &vr.ClassifyParams{}
		if ids := c.StringSlice("id"); len(ids) > 0 {
			params.ClassifierIDs = ids
		}
		if threshold := c.Float64("threshold"); threshold > 0.0 {
			params.Threshold = threshold
		}
		log.Printf("Requesting %d files for classifications...", len(sources))
		resp, err := client.ClassifyImages(context.Background(), xzip.NewArchiver(rawSources...), params)
		if err != nil {
			return err
		}
		if len(resp.Images) != len(sources) {
			return fmt.Errorf("unexpected response: watson returned no images without errors")
		}
		// rename image name as original
		for i := range resp.Images {
			resp.Images[i].Image = sources[i].orig
		}
		printImages(resp)
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
