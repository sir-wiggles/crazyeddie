package main

import (
	"archive/tar"
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html/charset"

	"github.com/yewno/MalfunctioningEddie/struct/arxiv"
	"github.com/yewno/MalfunctioningEddie/struct/brill"
	"github.com/yewno/MalfunctioningEddie/struct/nas"
	"github.com/yewno/publishers/acm"
)

const (
	NAS   = "nas"
	ACM   = "acm"
	ARXIV = "arxiv"
	BRILL = "brill"
)

type specific struct {
	Entry         []string
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)
}

var PublisherSpecifics = map[string]specific{
	NAS: specific{
		Entry: []string{"Product"},
	},
	ACM: specific{
		Entry:         []string{"proceeding", "periodical"},
		CharsetReader: charset.NewReaderLabel,
	},
	ARXIV: specific{
		Entry: []string{"arXiv"},
	},
	BRILL: specific{
		Entry: []string{"========"},
	},
}

func HandleZip(publisher string, reader *zip.Reader) error {

	for _, file := range reader.File {

	}
}

func HandleTAR(publisher string, reader *tar.Reader) error {

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			ext := path.Ext(header.Name)
			switch ext {

			case PDF:
				err = HandleTarPdf(publisher, reader, header.Name)

			case XML:
				temp, err := ioutil.TempFile("", "")
				if err != nil {
					return err
				}

				_, err = io.Copy(temp, reader)
				if err != nil {
					return err
				}

				decoder := xml.NewDecoder(temp)

				err = HandleXML(publisher, decoder)

			default:
				return errors.New(fmt.Sprintf("Need handler for ext (%s)", ext))
			}
		}
	}
	return nil
}

func HandleXML(publisher string, decoder *xml.Decoder) error {

	specs := PublisherSpecifics[publisher]

	if specs.CharsetReader != nil {
		decoder.CharsetReader = specs.CharsetReader
	}

	for {

		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch element := token.(type) {

		case xml.StartElement:

			log.Println(element.Name)
			if in(element.Name.Local, specs.Entry) {

				record := newRecord(publisher)
				err = decoder.DecodeElement(&record, &element)
				if err != nil {
					return err
				}

				document := newDocument(record)
				log.Println(document.Publisher, document.Record)
				time.Sleep(time.Second * 3)
			}
		}
	}
}

var arxivMetaURL = "http://export.arxiv.org/oai2?verb=GetRecord&identifier=oai:arXiv.org:%s&metadataPrefix=arXiv"
var arxividRegex = regexp.MustCompile(`((?i)[0-9a-z]{4,5}\.[0-9a-z]{4,5})`)

func HandleTarPdf(publisher string, file *tar.Reader, filename string) error {

	_, filename = path.Split(filename)

	id := strings.TrimSuffix(filename, ".pdf")
	if !arxividRegex.MatchString(id) {
		i := strings.IndexAny(id, "0123456789")
		id = id[0:i] + "/" + id[i:len(id)]
	}

	url := fmt.Sprintf(arxivMetaURL, id)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	return HandleXML(publisher, decoder)
}

func in(val string, items []string) bool {
	for _, item := range items {
		if val == item {
			return true
		}
	}
	return false
}

func newRecord(publisher string) interface{} {
	switch publisher {

	case NAS:
		return new(nas.Record)

	case ACM:
		return new(acm.Record)

	case ARXIV:
		return new(arxiv.Record)

	case BRILL:
		return new(brill.Record)
	}
	return nil
}

type document struct {
	Record    interface{}
	Publisher string
}

func newDocument(record interface{}) *document {
	switch record.(type) {

	case *nas.Record:
		return &document{
			Record:    record,
			Publisher: NAS,
		}

	case *acm.Record:
		return &document{
			Record:    record,
			Publisher: ACM,
		}

	case *arxiv.Record:
		return &document{
			Record:    record,
			Publisher: ARXIV,
		}
	}
	return nil
}
