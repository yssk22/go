package visualrecognition

import (
	"golang.org/x/net/context"
)

// FaceDetectResponse is a object returned by Datect* functions
type FaceDetectResponse struct {
	Images         []*Image `json:"images,omitempty"`
	ImageProcessed int      `json:"image_processed"`
}

// Image contains the result of face detection per an image
type Image struct {
	Faces       []*Face `json:"faces,omitempty"`
	Image       string  `json:"image,omitempty"`
	SourceURL   string  `json:"source_url,omitempty"`
	ResolvedURL string  `json:"resolved_url,omitempty"`
}

type Face struct {
	Age          *Age          `json:"age,omitempty"`
	FaceLocation *FaceLocation `json:"face_location,omitempty"`
	Gender       *Gender       `json:"gender,omitempty"`
	Identity     *Identity     `json:"identity,omitempty"`
}

type Age struct {
	Max   int     `json:"max"`
	Min   int     `json:"min"`
	Score float64 `json:"score"`
}

type FaceLocation struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
	Left   int64 `json:"left"`
	Top    int64 `json:"top"`
}

type Gender struct {
	Gender string  `json:"gender"`
	Score  float64 `json:"score"`
}

type Identity struct {
	Name          string  `json:"name"`
	Score         float64 `json:"score"`
	TypeHierarchy string  `json:"type_hierarchy"`
}

// DetectFacesByURL call /v3/detect_faces endpoint with the `url` parameter.
func (c *Client) DetectFacesByURL(ctx context.Context, url string) (*FaceDetectResponse, error) {
	var fdResp FaceDetectResponse
	if err := c.get("/v3/detect_faces", map[string][]string{
		"url": []string{url},
	}, &fdResp); err != nil {
		return nil, err
	}
	return &fdResp, nil
}
