package brill

type Record struct {
	ID     string `xml:"id,attr"`
	URL    string `xml:"targets,attr"`
	Volume string `xml:"volume,attr"`
	Page   string `xml:"page,attr"`
	Date   string `xml:"first-online,attr"`
	DOI    string `xml:"indo.doi,attr"`

	Title string `xml:"pseudoarticle>articleentry>mainentry"`
	Body  []struct {
		Text string `xml:",innerxml"`
	} `xml:"pseudoarticle>div"`
	Body2 []struct {
		Text string `xml:",innerxml"`
	} `xml:"pseudoarticle>p"`
	Body3 []struct {
		Text string `xml:",innerxml"`
	} `xml:"pseudoarticle>info"`
	//References  // listbibl
	// <pseudobiblio>
	//                 <bibliogroup>
	//                     <listbibl>
	//                         <p>Amma, K., <hi rend="italic">Mohiniyāṭṭam: Caritravum Ātaprakāravum</hi>, Kottayam, 1992 (Mal.).</p>
	//                     </listbibl>
	Authors []string `xml:"pseudoarticle>contributorgroup>name"`
}
