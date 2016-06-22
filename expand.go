package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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
)

var (
	// order matters here. Must check tar.gz before gz
	extensions = []string{ZIP, TAR_GZ, GZ, TAR}
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

	reader, err := decompress(filename)
	if err != nil {
		return err
	}

	switch reader.(type) {
	case *tar.Reader:

	case *gzip.Reader:

	case *zip.Reader:

	default:
		return errors.New(fmt.Sprintf("Expand got invalid type (%T)", reader))
	}
	err := Parse(publisher, nil)

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
