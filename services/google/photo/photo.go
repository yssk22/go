package photo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sync"

	"golang.org/x/oauth2"
)

// BaseURL is a base url for Google Photo API
const BaseURL = "https://picasaweb.google.com"

// RedirectURL is a string passed for oauth
const RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

// NewOAuth2Config returns a new *oauth2.Config for Google Photo endpoint.
func NewOAuth2Config(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		Scopes: []string{
			"https://picasaweb.google.com/data/",
		},
	}
}

// Client is a client for Gootle Photo API
type Client struct {
	BaseURL string
	client  *http.Client
}

// NewClient returns a new *Client
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		BaseURL: BaseURL,
		client:  client,
	}
}

// CreateAlbum creats a new album
func (c *Client) CreateAlbum(userID string, a *Album) (*Album, error) {
	if userID == "" {
		userID = "default"
	}
	postBody, _ := xml.Marshal(a)
	uri := fmt.Sprintf("%s/data/feed/api/user/%s", c.BaseURL, userID)
	resp, err := c.client.Post(uri, "application/atom+xml", bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error: %s", string(buff))
	}
	err = xml.NewDecoder(resp.Body).Decode(a)
	return a, err
}

// CreateAlbumByName creates a new album with the given album name
func (c *Client) CreateAlbumByName(userID string, albumName string) (*Album, error) {
	a := NewAlbum()
	a.Title = albumName
	return c.CreateAlbum(userID, a)
}

// ListAlbums returns a list of albums
func (c *Client) ListAlbums(userID string) ([]*Album, error) {
	if userID == "" {
		userID = "default"
	}
	uri := fmt.Sprintf("%s/data/feed/api/user/%s", c.BaseURL, userID)
	resp, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error: %s", string(buff))
	}
	return parseUserFeed(resp.Body)
}

// DeleteAlbum delets the album.
func (c *Client) DeleteAlbum(userID string, albumID string) error {
	if userID == "" {
		userID = "default"
	}
	uri := fmt.Sprintf("%s/data/entry/api/user/%s/albumid/%s", c.BaseURL, userID, albumID)
	req, _ := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("If-Match", "*")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API Error: %s", string(buff))
	}
	return nil
}

// ListMedia returns the list of media in an album
func (c *Client) ListMedia(userID, albumID string) ([]*Media, error) {
	if userID == "" {
		userID = "default"
	}
	uri := fmt.Sprintf("%s/data/feed/api/user/%s/albumid/%s", c.BaseURL, userID, albumID)
	resp, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error: %s", string(buff))
	}
	return parseAlbumFeed(resp.Body)
}

// GetMedia returns a *Media specified by mediaID
func (c *Client) GetMedia(userID, albumID, mediaID string) (*Media, error) {
	if userID == "" {
		userID = "default"
	}
	uri := fmt.Sprintf("%s/data/feed/api/user/%s/albumid/%s/photoid/%s", c.BaseURL, userID, albumID, mediaID)
	resp, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API Error: %s", string(buff))
	}
	return parsePhotoFeed(resp.Body)
}

// UploadMedia uploads a new media. Currently only photos are supported.
func (c *Client) UploadMedia(userID, albumID string, m *UploadMediaInfo, contentType string, contentLength int64, r io.Reader) (*Media, error) {
	if userID == "" {
		userID = "default"
	}
	if albumID == "" {
		albumID = "default"
	}
	const metaDataContentType = "application/atom+xml"
	var media *Media
	var err error
	var totalContentLength int64
	pipeOut, pipeIn := io.Pipe()
	uri := fmt.Sprintf("%s/data/feed/api/user/%s/albumid/%s", c.BaseURL, userID, albumID)
	writer := multipart.NewWriter(pipeIn)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	metaData, _ := xml.Marshal(m)
	totalContentLength += int64(3 + len(writer.Boundary()))        // --boundary\n
	totalContentLength += int64(14 + len(metaDataContentType) + 1) // Content-Type: xxxx\n\n
	totalContentLength += int64(len(metaData) + 1)                 // {xml}\n
	totalContentLength += int64(3 + len(writer.Boundary()))        // --boundary\n
	totalContentLength += int64(14 + len(contentType) + 1)         // Content-Type: xxxx\n\n
	totalContentLength += int64(len(metaData) + 1)                 // {binary}\n
	totalContentLength += int64(3 + len(writer.Boundary()))        // --boundary\n
	go func() {
		defer wg.Done()
		var resp *http.Response
		var req *http.Request
		req, err = http.NewRequest("POST", uri, pipeOut)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%q", writer.Boundary()))
		req.Header.Set("Content-Length", fmt.Sprintf("%d", totalContentLength))
		resp, err = c.client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 201 {
			buff, _ := ioutil.ReadAll(resp.Body)
			err = fmt.Errorf("API Error: %s", string(buff))
			return
		}
		media, err = parsePhotoFeed(resp.Body)
	}()

	metaDataPart, err := writer.CreatePart(textproto.MIMEHeader(map[string][]string{
		"Content-Type": []string{metaDataContentType},
	}))
	if err != nil {
		return nil, err
	}
	metaDataPart.Write(metaData)
	binaryPart, err := writer.CreatePart(textproto.MIMEHeader(map[string][]string{
		"Content-Type": []string{contentType},
	}))
	if _, err = io.Copy(binaryPart, r); err != nil {
		return nil, err
	}
	writer.Close()
	pipeIn.Close()
	wg.Wait()
	return media, err
}

// DeleteMedia delets a media
func (c *Client) DeleteMedia(userID string, albumID string, mediaID string) error {
	if userID == "" {
		userID = "default"
	}
	if albumID == "" {
		albumID = "default"
	}
	uri := fmt.Sprintf("%s/data/entry/api/user/%s/albumid/%s/photoid/%s", c.BaseURL, userID, albumID, mediaID)
	req, _ := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("If-Match", "*")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buff, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API Error: %s", string(buff))
	}
	return nil
}
