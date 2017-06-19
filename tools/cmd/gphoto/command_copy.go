// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/speedland/go/services/google/photo"
)

var commandCopy = cli.Command{
	Name:  "copy",
	Usage: "copy photos from URL",
	Action: func(c *cli.Context) {
		copy(c)
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "album, a",
			Usage: "album ID (default: \"default\")",
		},
		cli.StringFlag{
			Name:  "title, t",
			Usage: "Title (default: source url)",
		},
	},
}

func copy(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		printError(fmt.Errorf("No URLs are specified"))
	} else {
		var uploaders []*Copy
		albumID := c.String("album")
		for i := range args {
			title := c.String("title")
			if title == "" {
				title = args[i]
			}
			uploader, err := NewCopy(albumID, title, args[i])
			if err != nil {
				printError(err)
				continue
			}
			uploaders = append(uploaders, uploader)
		}
		bars := make([]*pb.ProgressBar, len(uploaders))
		for i, upload := range uploaders {
			bars[i] = upload.progressBar
		}
		p, _ := pb.StartPool(bars...)
		wg := new(sync.WaitGroup)
		for _, uploader := range uploaders {
			wg.Add(1)
			go func(u *Copy) {
				u.Execute()
				wg.Done()
			}(uploader)
		}
		wg.Wait()
		p.Stop()
		for _, uploader := range uploaders {
			if uploader.Media != nil {
				fmt.Printf("Media %s was uploaded.\n", uploader.Media.ID)
			} else {
				printError(fmt.Errorf("[%s] %v", uploader.Title, uploader.UploadError))
			}
		}
	}
}

type Copy struct {
	AlbumID       string
	Title         string
	URL           string
	ContentType   string
	ContentLength int64
	ByteRead      int64
	UploadError   error
	Media         *photo.Media
	client        *photo.Client
	response      *http.Response
	reader        io.Reader
	progressBar   *pb.ProgressBar
}

func (r *Copy) Execute() {
	defer r.progressBar.Finish()
	defer r.response.Body.Close()
	r.progressBar.Start()
	media, err := r.client.UploadMedia("", r.AlbumID, photo.NewUploadMediaInfo(r.Title), r.ContentType, r.ContentLength, r.reader)
	r.Media = media
	r.UploadError = err
}

func (r *Copy) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func NewCopy(albumID, title, url string) (*Copy, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s -> %s (Invalid status code)", url, resp.Status)
	}
	contentType := resp.Header.Get("content-type")
	contentSize := resp.ContentLength
	if !strings.HasPrefix(contentType, "image/") {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s -> %s (Invalid content type)", url, contentType)
	}

	progressBar := pb.New64(contentSize)
	progressBar.SetUnits(pb.U_BYTES)
	progressBar.Prefix(title)
	return &Copy{
		AlbumID:       albumID,
		Title:         title,
		URL:           url,
		ContentType:   contentType,
		ContentLength: contentSize,
		client:        NewClient(),
		response:      resp,
		reader:        progressBar.NewProxyReader(bufio.NewReaderSize(resp.Body, 4096*10)),
		progressBar:   progressBar,
	}, nil
}
