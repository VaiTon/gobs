package gobs

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func fromXml(res *http.Response, v any) (err error) {
	body := res.Body
	err = xml.NewDecoder(body).Decode(v)
	if err != nil {
		return fmt.Errorf("cannot decode xml: %w", err)
	}

	err = body.Close()
	if err != nil {
		return fmt.Errorf("cannot close body: %w", err)
	}

	return
}

func fromXmlToDir(res *http.Response) (dir Directory, err error) {
	err = fromXml(res, &dir)
	return
}
