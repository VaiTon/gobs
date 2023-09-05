package gobs

import (
	"fmt"
	"net/http"
)

var transport = http.DefaultTransport

type authenticatedTransport struct {
	username string
	password string
}

func (t *authenticatedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(t.username, t.password)
	res, err := transport.RoundTrip(request)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code not ok: %s", res.Status)
	}
	return res, err
}
