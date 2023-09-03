package gobs

import (
	"github.com/google/go-querystring/query"
	"io"
)

// GetBuild returns a simple listing of all repositories for the specified project.
//
// If project is an empty string, returns a listing of ALL repositories in ALL projects.
func (c *Client) GetBuild(proj string) ([]string, error) {
	res, err := c.http.Get(c.host + "/build/" + proj)
	if err != nil {
		return nil, err
	}

	dir, err := fromXmlToDir(res)
	if err != nil {
		return nil, err
	}

	repos := dir.GetEntryNames()
	return repos, nil
}

type BuildResultList struct {
	State   string        `xml:"state,attr"`
	Results []BuildResult `xml:"result"`
}

type BuildBinary struct {
	FileName string `xml:"filename,attr"`
	Size     string `xml:"size,attr"`
	MTime    string `xml:"mtime,attr"`
}

type BuildBinaryList struct {
	Package string        `xml:"package,attr"`
	Binary  []BuildBinary `xml:"binary"`
}

type BuildStatus struct {
	Package string `xml:"package,attr"`
	Code    string `xml:"code,attr"`
	Details string `xml:",chardata"`
}

type BuildSummary struct {
	StatusCount []struct {
		Code  string `xml:"code,attr"`
		Count string `xml:"count,attr"`
	} `xml:"statuscount"`
}

type BuildResult struct {
	Project    string            `xml:"project,attr"`
	Repository string            `xml:"repository,attr"`
	Arch       string            `xml:"arch,attr"`
	Code       string            `xml:"code,attr"`
	State      string            `xml:"state,attr"`
	Status     []BuildStatus     `xml:"status"`
	BinaryList []BuildBinaryList `xml:"binarylist"`
	Summary    BuildSummary      `xml:"summary"`
}

const (
	BuildResultQueryViewStatus     = "status"
	BuildResultQueryViewBinaryList = "binarylist"
	BuildResultQueryViewSummary    = "summary"
)

type BuildResultQuery struct {
	// View specifies which sections should be included in the result list.
	//  - [BuildResultQueryViewStatus]: Include detailed infos about the build status.
	//  - [BuildResultQueryViewSummary]: Include the summary of the status values.
	//  - [BuildResultQueryViewBinaryList]: Include a list of generated binary files.
	// If not specified the default value is [BuildResultQueryViewStatus].
	View       string `url:"view,omitempty"`
	Package    string `url:"package,omitempty"`
	Arch       string `url:"arch,omitempty"`
	Repository string `url:"repository,omitempty"`
	LastBuild  bool   `url:"lastbuild,omitempty"`
	LocalLink  bool   `url:"locallink,omitempty"`
	MultiBuild bool   `url:"multibuild,omitempty"`
}

func (c *Client) GetBuildResult(proj string, q BuildResultQuery) (*BuildResultList, error) {
	urlQuery, err := query.Values(q)
	if err != nil {
		return nil, err
	}

	encode := urlQuery.Encode()
	res, err := c.http.Get(c.host + "/build/" + proj + "/_result?" + encode)
	if err != nil {
		return nil, err
	}

	result := &BuildResultList{}
	if err := fromXml(res, result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetArchitectures returns a list of enabled architectures for a specific repository
// in a specific project
func (c *Client) GetArchitectures(proj string, repo string) ([]string, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo)
	if err != nil {
		return nil, err
	}

	dir, err := fromXmlToDir(res)
	if err != nil {
		return nil, err
	}

	archs := dir.GetEntryNames()
	return archs, nil
}

func (c *Client) GetPackageBinaries(proj, repo, arch, pkg string) (BuildBinaryList, error) {
	binaries := BuildBinaryList{Package: pkg}

	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg)
	if err != nil {
		return binaries, err
	}

	err = fromXml(res, &binaries)
	return binaries, err
}

func (c *Client) GetPackageBuildLog(proj, repo, arch, pkg string) (*string, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg + "/_log")
	if err != nil {
		return nil, err
	}

	log, err := io.ReadAll(res.Body)
	logStr := string(log)
	return &logStr, err
}

func (c *Client) GetPackageBuildStatus(proj, repo, arch, pkg string) (*BuildStatus, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg + "/_status")
	if err != nil {
		return nil, err
	}

	status := &BuildStatus{}
	err = fromXml(res, status)
	return status, err
}

type BuildReason struct {
	Explain   string `xml:"explain"`
	Time      string `xml:"time"`
	OldSource string `xml:"oldsource"`
}

func (c *Client) GetPackageBuildReason(proj, repo, arch, pkg string) (*BuildReason, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg + "/_reason")
	if err != nil {
		return nil, err
	}

	reason := &BuildReason{}
	err = fromXml(res, reason)
	return reason, err
}

type BuildJobStatus struct {
	Code         string `xml:"code,attr"`
	StartTime    string `xml:"starttime"`
	LastDuration string `xml:"lastduration"`
	HostArch     string `xml:"hostarch"`
	Uri          string `xml:"uri"`
	JobId        string `xml:"jobid"`
}

func (c *Client) GetPackageJobStatus(proj, repo, arch, pkg string) (*BuildJobStatus, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg + "/_jobstatus")
	if err != nil {
		return nil, err
	}

	status := &BuildJobStatus{}
	err = fromXml(res, status)
	return status, err
}

type BuildDimension struct {
	Unit  string `xml:"unit,attr"`
	Value string `xml:",chardata"`
}

func (b BuildDimension) String() string {
	return b.Value + " " + b.Unit
}

type BuildStatistics struct {
	Disk struct {
		Usage struct {
			Size       BuildDimension `xml:"size"`
			IoRequests string         `xml:"io_requests"`
			IoSectors  string         `xml:"io_sectors"`
		} `xml:"usage"`
	} `xml:"disk"`
	Memory struct {
		Usage struct {
			Size BuildDimension `xml:"size"`
		} `xml:"usage"`
	} `xml:"memory"`
	Times struct {
		Total struct {
			Time BuildDimension `xml:"time"`
		} `xml:"total"`
		Preinstall struct {
			Time BuildDimension `xml:"time"`
		} `xml:"preinstall"`
	} `xml:"times"`
	Download struct {
		Size      BuildDimension `xml:"size"`
		Binaries  int            `xml:"binaries"`
		CacheHits int            `xml:"cachehits"`
	} `xml:"download"`
}

func (c *Client) GetBuildStatistics(proj, repo, arch, pkg string) (*BuildStatistics, error) {
	res, err := c.http.Get(c.host + "/build/" + proj + "/" + repo + "/" + arch + "/" + pkg + "/_statistics")
	if err != nil {
		return nil, err
	}

	stats := &BuildStatistics{}
	err = fromXml(res, stats)
	return stats, err
}
