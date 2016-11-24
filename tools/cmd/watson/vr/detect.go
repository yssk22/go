package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/speedland/go/x/xarchive/xzip"
	"github.com/speedland/go/x/ximage"
	"github.com/urfave/cli"
)

var detect = cli.Command{
	Name:  "detect-faces",
	Usage: "deletct faces in a image files or URL",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "output",
			Usage: "output directory of detected faces",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) == 0 {
			return fmt.Errorf("must specify one url or a file")
		}
		if len(args) > maxFaceDetectionFiles {
			return fmt.Errorf("too many files or URLs")
		}
		output := c.String("output")
		tempdir, err := ioutil.TempDir("", "watson-vr")
		if err != nil {
			return fmt.Errorf("error creating a temp directory: %v", err)
		}
		defer os.RemoveAll(tempdir)
		rawSources, sources := prepareSources(tempdir, args...)
		if len(sources) == 0 {
			return fmt.Errorf("no files can be detected")
		}
		defer func() {
			for _, v := range rawSources {
				v.Close()
			}
		}()
		log.Printf("Requesting %d files for face detection...", len(sources))
		resp, err := client.DetectFacesOnImages(context.Background(), xzip.NewArchiver(rawSources...))
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
		printFaces(resp)
		if output != "" {
			return generateFaceImages(resp, sources, output)
		}
		return nil
	},
}

func printFaces(resp *vr.FaceDetectResponse) {
	for _, image := range resp.Images {
		fmt.Printf("%s:\n", path.Base(image.Image))
		for _, face := range image.Faces {
			if face.FaceLocation == nil {
				continue
			}
			fmt.Printf(
				"\t(%d, %d, %d, %d)\n",
				face.FaceLocation.Left, face.FaceLocation.Top, face.FaceLocation.Width, face.FaceLocation.Height,
			)
			if face.Age != nil {
				fmt.Printf("\t\tage: %d-%d (%f)\n", face.Age.Min, face.Age.Min, face.Age.Score)
			}
			if face.Gender != nil {
				fmt.Printf("\t\tgender: %s (%f)\n", face.Gender.Gender, face.Age.Score)
			}
			if face.Identity != nil {
				fmt.Printf("\t\tidentity: %s %s (%f)\n", face.Identity.Name, face.Identity.TypeHierarchy, face.Identity.Score)
			}
		}
	}
}

type sourceInfo struct {
	orig string
	temp string
}

func (s *sourceInfo) Path() string {
	if s.temp != "" {
		return s.temp
	}
	return s.orig
}

func generateFaceImages(resp *vr.FaceDetectResponse, sources []*sourceInfo, output string) error {
	if err := os.MkdirAll(output, os.FileMode(0755)); err != nil {
		return err
	}
	var hasSomeError = false
	for i, image := range resp.Images {
		sourcePath := sources[i].Path()
		sourceExt := path.Ext(sourcePath)
		baseName := strings.TrimSuffix(path.Base(sourcePath), sourceExt)
		for j, face := range image.Faces {
			faceFile := filepath.Join(output, fmt.Sprintf("%s-face-%d%s", baseName, j, sourceExt))
			fmt.Printf("Generating a cropped face file (%s)...", faceFile)
			err := func() error {
				src, err := os.Open(sourcePath)
				if err != nil {
					return err
				}
				defer src.Close()
				dst, err := os.Create(faceFile)
				if err != nil {
					return err
				}
				defer dst.Close()
				t := ximage.TypeByExtension(sourceExt)
				loc := face.FaceLocation
				return ximage.Crop(src, dst, t, int(loc.Left), int(loc.Top), int(loc.Left+loc.Width), int(loc.Top+loc.Height))
			}()
			if err != nil {
				fmt.Printf("error: %v\n", err)
				hasSomeError = true
			} else {
				fmt.Printf("done\n")
			}
		}
	}
	if hasSomeError {
		return fmt.Errorf("error generating some face files")
	}
	return nil
}

func download(url string, dir string) (string, error) {
	name := filepath.Join(dir, filepath.Base(url))
	file, err := os.Create(name)
	if err != nil {
		return "", fmt.Errorf("error download file: %v", err)
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("not 200 response")
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}
	return name, nil
}
