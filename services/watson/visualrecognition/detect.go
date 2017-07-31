package visualrecognition

import (
	"github.com/speedland/go/x/xarchive/xzip"
	"context"
)

// FaceDetectResponse is a object returned by Datect* functions.
type FaceDetectResponse struct {
	Images         []*FaceImage `json:"images,omitempty"`
	ImageProcessed int          `json:"image_processed"`
}

// FaceImage is detected faces info on a image.
type FaceImage struct {
	Faces       []*Face `json:"faces,omitempty"`
	Image       string  `json:"image,omitempty"`
	SourceURL   string  `json:"source_url,omitempty"`
	ResolvedURL string  `json:"resolved_url,omitempty"`
}

// Face is face info detected by watson.
type Face struct {
	Age          *Age          `json:"age,omitempty"`
	FaceLocation *FaceLocation `json:"face_location,omitempty"`
	Gender       *Gender       `json:"gender,omitempty"`
	Identity     *Identity     `json:"identity,omitempty"`
}

// Age is a detected age range with a score.
type Age struct {
	Max   int     `json:"max"`
	Min   int     `json:"min"`
	Score float64 `json:"score"`
}

// FaceLocation is a detected location of a face.
type FaceLocation struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
	Left   int64 `json:"left"`
	Top    int64 `json:"top"`
}

// Gender is a detected gender with a score.
type Gender struct {
	Gender string  `json:"gender"`
	Score  float64 `json:"score"`
}

// Identity is a detected identity with a score.
type Identity struct {
	Name          string  `json:"name"`
	Score         float64 `json:"score"`
	TypeHierarchy string  `json:"type_hierarchy"`
}

// DetectFacesOnURL make a GET request to `/v3/detect_faces`` with the `url` parameter.
func (c *Client) DetectFacesOnURL(ctx context.Context, url string) (*FaceDetectResponse, error) {
	var fdResp FaceDetectResponse
	if err := c.get("/v3/detect_faces", map[string][]string{
		"url": []string{url},
	}, &fdResp); err != nil {
		return nil, err
	}
	return &fdResp, nil
}

// DetectFacesOnImages makes a POST request to `/v3/detect_faces`` with a zipped image file given by sources.
func (c *Client) DetectFacesOnImages(ctx context.Context, source *xzip.Archiver) (*FaceDetectResponse, error) {
	var resp FaceDetectResponse
	files := make(map[string]*file)
	files["images_file"] = &file{
		name:   "image.zip",
		reader: source,
	}
	if err := c.postFiles("/v3/detect_faces", nil, nil, files, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
