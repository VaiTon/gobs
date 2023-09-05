package gobs

import "io"

type Package struct {
	Project string `xml:"project,attr"`
	Name    string `xml:"name,attr"`

	Issues []Issue `xml:"issue"`
}

type Issue struct {
	Change    string     `xml:"change,attr"`
	CreatedAt string     `xml:"created_at"`
	UpdatedAt string     `xml:"updated_at"`
	Name      string     `xml:"name"`
	Tracker   string     `xml:"tracker"`
	Label     string     `xml:"label"`
	Url       string     `xml:"url"`
	State     string     `xml:"state"`
	Summary   string     `xml:"summary"`
	Owner     IssueOwner `xml:"owner"`
}

type IssueOwner struct {
	RealName string `xml:"realname"`
	Email    string `xml:"email"`
	Login    string `xml:"login"`
}

func (c *Client) GetPackages(proj string) ([]string, error) {
	res, err := c.http.Get(c.host + "/source/" + proj)
	if err != nil {
		return nil, err
	}

	dir, err := fromXmlToDir(res)
	if err != nil {
		return nil, err
	}

	return dir.GetEntryNames(), nil
}

func (c *Client) GetSourceFile(project string, pkg string, file string) ([]byte, error) {

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
