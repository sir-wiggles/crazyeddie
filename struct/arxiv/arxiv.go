package arxiv

// http://arxiv.org/OAI/arXiv.xsd

type Record struct {
	ID string `xml:"id"`

	Created string `xml:"created"`
	Updated string `xml:"updated"`
	Authors []struct {
		FirstName   string   `xml:"forenames"`
		LastName    string   `xml:"keyname"`
		Affiliation []string `xml:"affiliation"`
	} `xml:"authors>author"`
	Title      string `xml:"title"`
	Categories string `xml:"categories"`
	DOI        string `xml:"doi"`
	Abstract   string `xml:"abstract"`
	License    string `xml:"license"`
}
