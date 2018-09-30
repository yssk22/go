package ximage

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestCrop(t *testing.T) {
	a := assert.New(t)
	var src bytes.Buffer
	var dst bytes.Buffer
	sample := image.NewRGBA(
		image.Rect(0, 0, 100, 100),
	)
	a.Nil(png.Encode(&src, sample))
	Crop(&src, &dst, ImageTypePNG, 25, 30, 35, 40)

	cropped, err := png.Decode(&dst)
	a.Nil(err)
	min := cropped.Bounds().Min
	max := cropped.Bounds().Max
	a.EqInt(0, min.X)
	a.EqInt(0, min.Y)
	a.EqInt(10, max.X)
	a.EqInt(10, max.Y)
}
