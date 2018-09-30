// Package ximage provider higher level functions for image processing
package ximage

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
)

// ImageType is a type string of an image
type ImageType string

// Available ImageType values
const (
	ImageTypeUnknown ImageType = "unknown"
	ImageTypePNG     ImageType = "png"
	ImageTypeJPEG    ImageType = "jpeg"
)

// TypeByExtension returns ImageType by extention
func TypeByExtension(ext string) ImageType {
	switch mime.TypeByExtension(ext) {
	case "image/jpeg":
		return ImageTypeJPEG
	case "image/png":
		return ImageTypePNG
	}
	return ImageTypeUnknown
}

// Decode decodes the image
func Decode(src io.Reader, t ImageType) (image.Image, error) {
	switch t {
	case ImageTypePNG:
		return png.Decode(src)
	case ImageTypeJPEG:
		return jpeg.Decode(src)
	default:
		return nil, ErrUnsupportedImageType
	}
}

// ErrUnsupportedImageType is an error when unsupported ImageType passed
var ErrUnsupportedImageType = fmt.Errorf("unsupported image type")
