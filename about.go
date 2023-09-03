package gobs

type About struct {
	Title          string `xml:"title"`
	Description    string `xml:"description"`
	Revision       string `xml:"revision"`
	LastDeployment string `xml:"last_deployment"`
	Commit         string `xml:"commit"`
}

func (c *Client) GetAbout() (about About, err error) {
	res, err := c.http.Get(c.host + "/about")
	if err != nil {
		return
	}
	err = fromXml(res, &about)
	if err != nil {
		return
	}
	return
}
