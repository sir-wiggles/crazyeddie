package acm

type Record struct {
	// Periodical
	Journal ACMJournal `xml:"journal_rec"`
	Issue   ACMIssue   `xml:"issue_rec"`

	// Proceeding
	Conference ACMConference `xml:"conference_rec"`
	Proceeding ACMProceed    `xml:"proceeding_rec"`

	// Common to both
	Sections []ACMSection `xml:"content>section"`
	Articles []ACMArticle `xml:"content>article_rec"`
}

type ACMArticle struct {
	ID           string `xml:"article_id"`
	Title        string `xml:"title"`
	SubTitle     string `xml:"subtitle"`
	Date         string `xml:"article_publication_date"`
	StartPage    string `xml:"page_from"`
	EndPage      string `xml:"page_to"`
	ManuscriptId string `xml:"manuscript_tracking_id"`
	DOI          string `xml:"doi_number"`
	URL          string `xml:"url"`
	Abstract     struct {
		Text string `xml:",innerxml"`
	} `xml:"abstract>par"`
	Keywords []struct {
		Text string `xml:",innerxml"`
	} `xml:"keywords>kw"`
	// Categories struct {
	// 	// Primary
	// 	//Other
	// 	}`xml:"categories"`
	Authors []struct {
		ID         string `xml:"person_id"`
		Seq        string `xml:"seq_no"`
		FirstName  string `xml:"first_name"`
		LastName   string `xml:"last_name"`
		MiddleName string `xml:"middle_name"`
		Suffix     string `xml:"suffix"`
		Role       string `xml:"role"`
	} `xml:"authors>au"`
	Body []struct {
		Text string `xml:",innerxml"`
	} `xml:"fulltext>ft_body"`
}

type ACMPublisher struct {
	ID   string `xml:"publisher_id"`
	Code string `xml:"publisher_code"`
	Name string `xml:"publisher_name"`
	URL  string `xml:"publisher_url"`
}

type ACMSection struct {
	ID    string `xml:"section_id"`
	Title string `xml:"section_title"`
	Type  string `xml:"section_type"`
	//chair_editor
	Articles []ACMArticle `xml:"article_rec"`
}
type ACMJournal struct {
	ID        string       `xml:"journal_id"`
	Code      string       `xml:"journal_code"`
	Name      string       `xml:"journal_name"`
	Abbr      string       `xml:"journal_abbr"`
	ISSN      string       `xml:"issn"`
	EISSN     string       `xml:"eissn"`
	Language  string       `xml:"language"`
	Type      string       `xml:"periodical_type"`
	Publisher ACMPublisher `xml:"publisher"`
}

type ACMIssue struct {
	ID     string `xml:"issue_id"`
	Volume string `xml:"volume"`
	Issue  string `xml:"issue"`
	Date   string `xml:"publication_date"`
}

// periodical.dtd
type ACMPeriodical struct {
	Journal   ACMJournal `xml:"journal_rec"`
	Issue     ACMIssue   `xml:"issue_rec"`
	SectionID string     `xml:"section_id"`
	Title     string     `xml:"section_title"`
	Type      string     `xml:"section_type"`
	Article   ACMArticle `xml:"article_rec"`
}

type ACMConference struct {
	StartDate string `xml:"conference_date>start_date"`
	EndDate   string `xml:"conference_date>end_date"`
	URL       string `xml:"conference_url"`
	Location  struct {
		City    string `xml:"city"`
		State   string `xml:"state"`
		Country string `xml:"country"`
	} `xml:"conference_loc"`
}

type ACMProceed struct {
	ISBN        string `xml:"isbn"`
	ISBN13      string `xml:"isbn13"`
	ISSN        string `xml:"issn"`
	EISSN       string `xml:"eissn"`
	Date        string `xml:"publication_date"`
	Title       string `xml:"proc_title"`
	Description string `xml:"proc_desc"`
	Abstract    struct {
		Text string `xml:",innerxml"`
	} `xml:"abstract>par"`
	Publisher ACMPublisher `xml:"publisher"`
	Copyright struct {
		Holder string `xml:"copyright_holder_name"`
		Year   string `xml:"copyright_holder_year"`
	} `xml:"ccc>copyright_holder"`
}

type ACMProceeding struct {
	Conference ACMConference `xml:"conference_rec"`
	Proceeding ACMProceed    `xml:"proceeding_rec"`
	SectionID  string        `xml:"section_id"`
	Title      string        `xml:"section_title"`
	Type       string        `xml:"section_type"`

	Article ACMArticle `xml:"article_rec"`
}
