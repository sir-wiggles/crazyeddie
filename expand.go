package main

import (
	"encoding/xml"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var (
		publisher string = "nas"
		filename  string = "nas.xml"
	)

	if err := Expand(publisher, filename); err != nil {
		log.Println(err)
	}
}

func Expand(publisher, filename string) error {

	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(f)

	err = Parse(publisher, decoder)

	return nil
}
