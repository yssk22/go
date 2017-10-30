package facebook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/speedland/go/x/xerrors"
)

type SlideshowParams struct {
	ImagesURLs   []string `json:"images_urls"`
	DurationMs   int      `json:"duration_ms,omitempty"`
	TransitionMs int      `json:"transition_ms,omitempty"`
}

func NewSlideshowParams(imagesURLs []string) *SlideshowParams {
	return &SlideshowParams{
		ImagesURLs:   imagesURLs,
		DurationMs:   1750,
		TransitionMs: 250,
	}
}

// CreateSlideshow creates a slideshow video on the page specified by page id.
func (c *Client) CreateSlideshow(ctx context.Context, id string, params *SlideshowParams) (string, error) {
	var r map[string]string
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	specJSON, err := json.Marshal(params)
	if err != nil {
		w.Close()
		return "", xerrors.Wrap(err, "could not marshal *SlideshowParams")
	}
	if err := w.WriteField("slideshow_spec", string(specJSON)); err != nil {
		w.Close()
		return "", xerrors.Wrap(err, "could not write multipart field `slideshow_spec`")
	}
	w.Close()
	if err := c.PostVideo(ctx, fmt.Sprintf("/%s/videos", id), nil, &MultipartParams{
		ContentType: w.FormDataContentType(),
		Body:        &body,
	}, &r); err != nil {
		return "", err
	}
	return r["id"], nil
}
