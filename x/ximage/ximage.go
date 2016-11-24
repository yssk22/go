// Package ximage provider higher level functions for image processing
package ximage

import (
	"fmt"
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

// ErrUnsupportedImageType is an error when unsupported ImageType passed
var ErrUnsupportedImageType = fmt.Errorf("unsupported image type")
