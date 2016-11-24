// watson-vr is a cli command to manage watson visual recognition
//
// - watson-vr classify {url}: classify a image on {url}.
// - watson-vr detect-face {url}: detect faces in a image on url
// - watson-vr classifiers list: list all custom classifiers
// - watson-vr classifiers create {name} {directory}: create a new custom classifier
// - watson-vr classifiers show {classifier-id}: show a custom classifier
// - watson-vr classifiers delete {classifier-id}: delete a custom classifier
//

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/speedland/go/x/xtime"
	"github.com/urfave/cli"
)

var apiKey string

var Version = "master"
var BuildRev = "unknown"

func main() {
	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()
	log.SetPrefix("")
	log.SetFlags(0)
	log.SetOutput(devNull)

	app := cli.NewApp()
	app.Name = "watson-vr"
	app.Usage = "manage watson visual recognition"
	app.Version = fmt.Sprintf("%s-%s", Version, BuildRev)
	app.Commands = []cli.Command{
		classify,
		detect,
		prepareClassifier,
		createClassifier,
		listClassifiers,
		showClassifier,
		updateClassifier,
		deleteClassifier,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apikey, k",
			Usage:  "API key to access watson",
			EnvVar: "WATSON_API_KEY",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose logging",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("verbose") {
			log.SetOutput(os.Stderr)
		}
		return nil
	}

	app.Run(os.Args)
}

func NewClient(c *cli.Context) (*vr.Client, error) {
	apikey := c.GlobalString("apikey")
	if apikey == "" {
		return nil, fmt.Errorf("--apikey must be speicified")
	}
	return vr.NewClient(apikey, http.DefaultClient), nil
}

var userLocation = time.Now().Location()

func formatTime(t time.Time) string {
	return xtime.FormatDateTimeString(t.In(userLocation))
}
