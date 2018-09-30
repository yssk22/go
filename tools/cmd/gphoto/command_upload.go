// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"path/filepath"

	"github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/yssk22/go/services/google/photo"
)

var commandUpload = cli.Command{
	Name:  "upload",
	Usage: "upload photos",
	Action: func(c *cli.Context) {
		upload(c)
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "album, a",
			Usage: "album ID (default: \"default\")",
		},
	},
}

func upload(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		printError(fmt.Errorf("No files are specified"))
	} else {
		albumID := c.String("album")
		uploaders := make([]*Uploader, 0)
		for i, _ := range args {
			title := filepath.Base(args[i])
			uploader, err := NewUploader(albumID, title, args[i])
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
			go func(u *Uploader) {
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

type Uploader struct {
	AlbumID       string
	Title         string
	FilePath      string
	ContentType   string
	ContentLength int64
	ByteRead      int64
	UploadError   error
	Media         *photo.Media
	client        *photo.Client
	file          *os.File
	reader        io.Reader
	progressBar   *pb.ProgressBar
}

func (r *Uploader) Execute() {
	defer r.progressBar.Finish()
	defer r.file.Close()
	r.progressBar.Start()
	media, err := r.client.UploadMedia("", r.AlbumID, photo.NewUploadMediaInfo(r.Title), r.ContentType, r.ContentLength, r.reader)
	r.Media = media
	r.UploadError = err
}

func (r *Uploader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func NewUploader(albumID, title, path string) (*Uploader, error) {
	contentType, err := getContentType(path)
	if err != nil {
		return nil, err
	}
	contentSize, err := getContentSize(path)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	progressBar := pb.New64(contentSize)
	progressBar.SetUnits(pb.U_BYTES)
	progressBar.Prefix(title)
	return &Uploader{
		AlbumID:       albumID,
		Title:         title,
		FilePath:      path,
		ContentType:   contentType,
		ContentLength: contentSize,
		client:        NewClient(),
		file:          file,
		reader:        progressBar.NewProxyReader(bufio.NewReaderSize(file, 4096*10)),
		progressBar:   progressBar,
	}, nil
}

var contentTypeMapper = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".ts":   "video/mpeg",
	".mp4":  "video/mpeg4",
	".mpeg": "video/mpeg",
	".3gp":  "video/3gpp",
	".avi":  "video/avi",
}

func getContentType(file string) (string, error) {
	ext := filepath.Ext(file)
	if mime, ok := contentTypeMapper[ext]; ok {
		return mime, nil
	} else {
		return "", fmt.Errorf("Unsupported file type.")
	}
}

func getContentSize(file string) (int64, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	if stat.IsDir() {
		return 0, fmt.Errorf("Not a file")
	}
	return stat.Size(), nil
}
