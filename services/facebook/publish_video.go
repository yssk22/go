package facebook

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/yssk22/go/x/xtime"

	"github.com/yssk22/go/x/xerrors"
)

type VideoParams struct {
	Title                string           `json:"title,omitempty"`
	Description          string           `json:"description,omitempty"`
	FileURL              string           `json:"file_url,omitempty"`
	Published            *bool            `json:"published,omitempty"`
	ScheduledPublishTime *xtime.Timestamp `json:"scheduled_publish_time,omitempty"`
	SponsorID            string           `json:"sponsor_id,omitempty"`
	SlideshowSpec        *SlideshowSpec   `json:"slideshow_spec,omitempty"`
}

type SlideshowSpec struct {
	ImagesURLs   []string `json:"images_urls"`
	DurationMs   int      `json:"duration_ms,omitempty"`
	TransitionMs int      `json:"transition_ms,omitempty"`
}

func NewSlideshowParams(imagesURLs []string) *VideoParams {
	return &VideoParams{
		SlideshowSpec: &SlideshowSpec{
			ImagesURLs:   imagesURLs,
			DurationMs:   1750,
			TransitionMs: 250,
		},
	}
}

// CreateSlideshow creates a slideshow video on the page specified by page id.
func (c *Client) CreateSlideshow(ctx context.Context, id string, params *VideoParams) (string, error) {
	var r map[string]string
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := writeStructIntoMultipart(w, params); err != nil {
		w.Close()
		return "", xerrors.Wrap(err, "could not write multipart field of *VideoParams")
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
