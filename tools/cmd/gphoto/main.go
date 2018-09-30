// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"log"
	"os"

	"io/ioutil"

	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/yssk22/go/x/xos"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "yssk22-gphoto"
	app.Usage = "Google Photo Client for yssk22 apps"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "./.gphotorc",
			Usage:  "configuration file",
			EnvVar: "CONFIG",
		},
	}
	app.Before = func(c *cli.Context) error {
		paths := []string{c.String("config"), "~/.gphotorc", "/etc/gphotorc"}
		path, err := xos.TryExists(paths...)
		if err != nil {
			log.Fatalf("None of %s exist", paths)
		}
		buff, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if err = json.Unmarshal(buff, &DefaultConfig); err != nil {
			log.Fatalf(err.Error())
		}
		if DefaultConfig.ClientID == "" ||
			DefaultConfig.ClientSecret == "" ||
			DefaultConfig.AccessToken == "" {
			log.Println("One of client_id, client_secret, access_token is missing in %s", path)
			log.Println("Parsed content: ")
			log.Fatalln(string(buff))
		}
		return nil
	}
	app.Commands = []cli.Command{
		commandLs,
		commandMkalbum,
		commandUpload,
		commandCopy,
		commandRm,
	}
	app.Run(os.Args)
}
