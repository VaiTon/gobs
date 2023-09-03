package gobs

type Directory struct {
	Count   int16   `xml:"count,attr"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Name string `xml:"name,attr"`
}

func (dir Directory) GetEntryNames() (names []string) {
	names = make([]string, 0, dir.Count)
	for _, e := range dir.Entries {
		names = append(names, e.Name)
	}
	return
}
