package gobs

import "github.com/google/go-querystring/query"

func (c *Client) GetRequests() ([]string, error) {

	res, err := c.http.Get(c.host + "/request")
	if err != nil {
		return nil, err
	}

	dir, err := fromXmlToDir(res)
	if err != nil {
		return nil, err
	}

	return dir.GetEntryNames(), nil
}

type RequestCollection struct {
	Requests []Request `xml:"request"`
	Matches  int       `xml:"matches,attr"`
}

type PackageRef struct {
	Package string `xml:"package,attr"`
	Project string `xml:"project,attr"`
	Rev     string `xml:"rev,attr"`
}

type RequestHistory struct {
	Who         string `xml:"who,attr"`
	When        string `xml:"when,attr"`
	Description string `xml:"description"`
	Comment     string `xml:"comment"`
}

type RequestAction struct {
	Type   string     `xml:"type,attr"`
	Source PackageRef `xml:"source"`
	Target PackageRef `xml:"target"`
}

type RequestState struct {
	Name    string `xml:"name,attr"`
	Who     string `xml:"who,attr"`
	When    string `xml:"when,attr"`
	Comment string `xml:",chardata"`
}

type Request struct {
	Id          string           `xml:"id,attr"`
	Creator     string           `xml:"creator,attr"`
	Action      []RequestAction  `xml:"action"`
	State       RequestState     `xml:"state"`
	Review      []Review         `xml:"review"`
	History     []RequestHistory `xml:"history"`
	Description string           `xml:"description"`
}

type Review struct {
	State  string `xml:"state,attr"`
	Who    string `xml:"who,attr"`
	When   string `xml:"when,attr"`
	ByUser string `xml:"by_user,attr"`
}

func (c *Client) GetRequest(id string) (*Request, error) {

	res, err := c.http.Get(c.host + "/request/" + id)
	if err != nil {
		return nil, err
	}

	req := &Request{}
	err = fromXml(res, req)

	return req, err
}

type RequestCollectionQuery struct {
	User            string   `url:"user,omitempty"`
	Project         string   `url:"project,omitempty"`
	Package         string   `url:"package,omitempty"`
	States          []string `url:"states,omitempty" del:","`
	Types           []string `url:"types,omitempty" del:","`
	Roles           []string `url:"roles,omitempty" del:","`
	WithHistory     bool     `url:"withhistory,int,omitempty"`
	WithFullHistory bool     `url:"withfullhistory,int,omitempty"`
	Limit           int      `url:"limit,omitempty"`
	Offset          int      `url:"offset,omitempty"`
	Ids             []string `url:"ids,omitempty" del:","`
}

func (c *Client) GetRequestCollection(q RequestCollectionQuery) (*RequestCollection, error) {
	values, err := query.Values(q)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Get(c.host + "/request?view=collection&" + values.Encode())
	if err != nil {
		return nil, err
	}

	collection := &RequestCollection{}
	err = fromXml(res, collection)
	return collection, err
}
