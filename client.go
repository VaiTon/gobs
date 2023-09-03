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

func (c *Client) GetFile(project string, pkg string, file string) ([]byte, error) {

	res, err := c.http.Get(c.host + "/source/" + project + "/" + pkg + "/" + file)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
