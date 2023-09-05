package gobs

import (
	"io"
	"net/http"
)

type Client struct {
	http *http.Client
	host string
}

func NewClient(host string, username string, password string) *Client {
	t := &authenticatedTransport{username: username, password: password}
	return &Client{host: host, http: &http.Client{Transport: t}}
}

// GetRaw returns the raw response body for a given URL
//
// This is useful for downloading files or other content without parsing it into
// the corresponding struct.
//
// The request is authenticated using the credentials provided when creating the
// client.
func (c *Client) GetRaw(url string) (io.ReadCloser, error) {
	res, err := c.http.Get(c.host + url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
