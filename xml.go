package main

import (
	"encoding/xml"
	"io"
	"log"
	"time"
)

var PublisherMappings = map[string][]string{
	"nas": {"Product"},
	"acm": {"proceeding", "periodical"},
}

func Parse(publisher string, decoder *xml.Decoder) error {

	mappings := PublisherMappings[publisher]

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

			if in(element.Name.Local, mappings) {

				record := newRecord(publisher)
				err = decoder.DecodeElement(&record, &element)
				if err != nil {
					return err
				}

				log.Println(string(record.(*nasRecord).Content))
				time.Sleep(time.Second * 1)
			}
		}
	}
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
	case "nas":
		return new(nasRecord)
	}
	return nil
}

func newDocument(record interface{}) Documenter {
	return nil
}

type Documenter interface {
	Content() []byte
}
