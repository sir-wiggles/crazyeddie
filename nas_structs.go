package main

import "encoding/xml"

type nasRecord struct {
	XMLName  xml.Name `xml:"Product"`
	Content  []byte   `xml:",innerxml"`
	YewnoIDs []struct {
		Type  string `xml:"ProductIDType"`
		Value string `xml:"IDValue"`
	} `xml:"ProductIdentifier"`
}
