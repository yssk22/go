// Package xzip provides higher level types and functions on top of "arvhive/zip" package.
package xzip

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xtime"
)

const uint32max = (1 << 32) - 1

// Archiver is a object to build a zip. You can read multiple source stream and forward as a zip stream
// to the destination like file, http form, ...etc. Builder is implemented on top of io.Pile so you do
// zip compression with less buffering on memory.
//
// The basic usage is like:
//
//     a, _ := NewRawSourceFromFile("a.txt")
//     defer a.Close()
//     b, _ := NewRawSourceFromFile("b.txt")
//     defer b.Close()
//
//     builder := NewArchiver(a, b)
//     zip, _ := os.Create("ab.zip")
//     defer zip.Close()
//     io.Copy(zip, buiilder)
//
type Archiver struct {
	sources    []*RawSource
	errors     []error
	pipeReader io.Reader
}

// Read reads the zip content from source.
func (s *Archiver) Read(p []byte) (int, error) {
	n, err := s.pipeReader.Read(p)
	if err != nil {
		return n, err
	}
	if len(s.errors) != 0 {
		return n, fmt.Errorf("zip source is broken: %v", s.errors[0])
	}
	return n, nil
}

// Close closes underlying sources
func (s *Archiver) Close() error {
	me := xerrors.NewMultiError(len(s.sources))
	for i, s := range s.sources {
		me[i] = s.Close()
	}
	if me.HasError() {
		return me
	}
	return nil
}

// NewArchiver returns a new *Archiver from multiple raw sources.
func NewArchiver(sources ...*RawSource) *Archiver {
	archiver := &Archiver{
		sources: sources,
		errors:  make([]error, 0),
	}
	pipeReader, pipeWriter := io.Pipe()
	archiver.pipeReader = pipeReader
	go func() {
		zipWriter := zip.NewWriter(pipeWriter)
		defer pipeWriter.Close()
		defer zipWriter.Close()
		for _, source := range archiver.sources {
			header := source.ZipHeader()
			header.Method = zip.Deflate
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				archiver.errors = append(archiver.errors, err)
				continue
			}
			if _, err := io.Copy(writer, source); err != nil {
				archiver.errors = append(archiver.errors, err)
				continue
			}
		}
	}()
	return archiver
}

// RawSource is a source to create a Zip stream.
type RawSource struct {
	Name       string
	Size       uint64
	ModifiedAt time.Time
	Mode       os.FileMode
	source     io.ReadCloser
}

// Read reads the data from stream
func (s *RawSource) Read(p []byte) (int, error) {
	return s.source.Read(p)
}

// Close close the source stream
func (s *RawSource) Close() error {
	return s.source.Close()
}

// ZipHeader returns *zip.FileHeader to create a zip stream
func (s *RawSource) ZipHeader() *zip.FileHeader {
	fh := &zip.FileHeader{
		Name:               s.Name,
		UncompressedSize64: s.Size,
	}
	fh.SetModTime(s.ModifiedAt)
	fh.SetMode(s.Mode)
	if s.Size > uint32max {
		fh.UncompressedSize = uint32max
	} else {
		fh.UncompressedSize = uint32(fh.UncompressedSize64)
	}
	return fh
}

// NewRawSourceFromFile returns a new *RawSource from file path
func NewRawSourceFromFile(path string) (*RawSource, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is not a file", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	source := &RawSource{
		Name:       info.Name(),
		Size:       uint64(info.Size()),
		ModifiedAt: info.ModTime(),
		Mode:       info.Mode(),
		source:     file,
	}
	return source, nil
}

// NewRawSourceFromURL returns a new *RawSource from URL
func NewRawSourceFromURL(url string, client *http.Client) (*RawSource, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response (%d)", resp.StatusCode)
	}
	return &RawSource{
		Name:       path.Base(url),
		Size:       uint64(resp.ContentLength),
		Mode:       os.FileMode(0644),
		ModifiedAt: xtime.Now(),
		source:     resp.Body,
	}, nil
}
