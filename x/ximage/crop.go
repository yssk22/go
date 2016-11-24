package ximage

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
)

// Crop crops the source image and write the cropped image to dst.
func Crop(src io.Reader, dst io.Writer, t ImageType, x0, y0, x1, y1 int) error {
	var srcImg image.Image
	var err error
	switch t {
	case ImageTypePNG:
		srcImg, err = png.Decode(src)
		break
	case ImageTypeJPEG:
		srcImg, err = jpeg.Decode(src)
		break
	default:
		return ErrUnsupportedImageType
	}
	if err != nil {
		return err
	}
	crop := srcImg.Bounds().Intersect(image.Rect(x0, y0, x1, y1))
	result := image.NewRGBA(crop)
	draw.Draw(result, crop, srcImg, crop.Min, draw.Src)
	switch t {
	case ImageTypePNG:
		return png.Encode(dst, result)
	case ImageTypeJPEG:
		return jpeg.Encode(dst, result, &jpeg.Options{
			Quality: 100,
		})
	default:
		return ErrUnsupportedImageType
	}
}
