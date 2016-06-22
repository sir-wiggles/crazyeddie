package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ZIP    string = ".zip"
	TAR_GZ string = ".tar.gz"
	GZ     string = ".gz"
	TAR    string = ".tar"
	XML    string = ".xml"
	PDF    string = ".pdf"
)

var (
	// order matters here. Must check tar.gz before gz
	extensions = []string{ZIP, TAR_GZ, GZ, TAR, XML}
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var (
		publisher string = "arxiv"
		filename  string = fmt.Sprintf("./struct/%s/%s.tar", publisher, publisher)
	)

	if err := Expand(publisher, filename); err != nil {
		log.Println(err)
	}
}

func Expand(publisher, filename string) error {

	reader, err := decompress(filename)
	if err != nil {
		return err
	}

	switch reader.(type) {
	case *tar.Reader:
		err = HandleTAR(publisher, reader.(*tar.Reader))

	case *gzip.Reader:
		decoder := xml.NewDecoder(reader.(*gzip.Reader))
		err = HandleXML(publisher, decoder)

	case *zip.Reader:
		log.Println("zip case not implemented")

	case *os.File:
		decoder := xml.NewDecoder(reader.(*os.File))
		err = HandleXML(publisher, decoder)

	default:
		return errors.New(fmt.Sprintf("Expand got invalid type (%T)", reader))
	}

	return err
}

func decompress(filename string) (interface{}, error) {

	var (
		reader interface{}
		err    error
	)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	for _, ext := range extensions {

		if !strings.HasSuffix(filename, ext) {
			continue
		}

		switch ext {

		case ZIP:
			reader, err = unzip(file)

		case GZ:
			reader, err = ungzip(file)

		case TAR:
			reader, err = untar(file)

		case TAR_GZ:
			reader, err = untargz(file)

		case XML:
			return file, nil
		}

		if err != nil {
			return nil, err
		}
		return reader, err
	}
	return nil, errors.New(fmt.Sprintf("invalid extension for filename (%s)", filename))
}

func unzip(file *os.File) (*zip.Reader, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func ungzip(file *os.File) (*gzip.Reader, error) {
	return gzip.NewReader(file)
}

func untar(file *os.File) (*tar.Reader, error) {
	return tar.NewReader(file), nil
}

func untargz(file *os.File) (*tar.Reader, error) {
	reader, err := ungzip(file)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(reader), nil
}
