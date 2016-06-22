package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/satori/go.uuid"
)

var decompressTest = map[string]func(*testing.T) string{
	"zip":    zipFile,
	"gzip":   gzipFile,
	"tar":    tarFile,
	"tar.gz": tarGzFile,
	"xml":    normal,
}

func TestDecompress(t *testing.T) {

	for _, v := range decompressTest {
		filename := v(t)
		reader, _ := decompress(filename)

		switch reader.(type) {
		case *gzip.Reader:
			r := reader.(*gzip.Reader)
			buf := bytes.NewBuffer(nil)
			io.Copy(buf, r)
			log.Println(buf.String())

		case *zip.Reader:
			r := reader.(*zip.Reader)
			for _, f := range r.File {
				log.Println(f)
			}

		case *tar.Reader:
			r := reader.(*tar.Reader)
			for {
				h, err := r.Next()
				if err != nil {
					if err == io.EOF {
						break
					}
				}
				log.Println(h)
			}
		case *os.File:
			r := reader.(*os.File)
			log.Println(r)

		default:
			log.Printf("%T", reader)
		}
	}
}

func normalFile(t *testing.T) string {
	file, err := tempFile("xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	file.Write([]byte("normal file"))
	return file.Name()
}

func tarFile(t *testing.T) string {
	file, err := tempFile("tar")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := _tar(t)

	_, err = io.Copy(file, reader)
	if err != nil {
		t.Fatal(err)
	}

	return file.Name()
}

func tarGzFile(t *testing.T) string {

	file, err := tempFile("tar.gz")
	if err != nil {
		t.Fatal(err)
	}

	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)

	var files = []struct {
		Name, Body string
	}{
		{"foo.txt", "foo bar baz."},
		{"bar.txt", "ni ni ni"},
		{"baz.txt", "green eggs and ham"},
	}
	for _, file := range files {

		h := &tar.Header{
			Name:  file.Name,
			Size:  int64(len(file.Body)),
			Uname: os.Getenv("USER"),
			Gname: os.Getenv("USER"),
			Mode:  664,
		}

		err := tw.WriteHeader(h)
		if err != nil {
			t.Fatal(err)
		}

		_, err = tw.Write([]byte(file.Body))
		if err != nil {
			t.Fatal(err)
		}

	}
	err = tw.Close()
	if err != nil {
		t.Fatal(err)
	}
	return file.Name()
}

func gzipFile(t *testing.T) string {
	file, err := tempFile("gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := _gz(t)

	_, err = io.Copy(file, reader)
	if err != nil {
		t.Fatal(err)
	}

	return file.Name()
}

func zipFile(t *testing.T) string {

	file, err := tempFile("zip")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := _zip(t)

	_, err = io.Copy(file, reader)
	if err != nil {
		t.Fatal(err)
	}

	return file.Name()
}

func _zip(t *testing.T) io.Reader {

	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"foo.txt", "foo bar baz."},
		{"bar.txt", "ni ni ni"},
		{"baz.txt", "green eggs and ham"},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewReader(buf.Bytes())
}

func _gz(t *testing.T) io.Reader {

	buf := new(bytes.Buffer)

	w := gzip.NewWriter(buf)

	_, err := w.Write([]byte("foo bar baz"))
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(buf.Bytes())
}

func _tar(t *testing.T) io.Reader {

	buf := new(bytes.Buffer)

	w := tar.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"foo.txt", "foo bar baz."},
		{"bar.txt", "ni ni ni"},
		{"baz.txt", "green eggs and ham"},
	}
	for _, file := range files {

		h := &tar.Header{
			Name:  file.Name,
			Size:  int64(len(file.Body)),
			Uname: os.Getenv("USER"),
			Gname: os.Getenv("USER"),
			Mode:  664,
		}

		err := w.WriteHeader(h)
		if err != nil {
			t.Fatal(err)
		}

		_, err = w.Write([]byte(file.Body))
		if err != nil {
			t.Fatal(err)
		}

	}
	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewReader(buf.Bytes())
}

func tempFile(ext string) (*os.File, error) {
	fn := uuid.NewV4()
	return os.Create(fmt.Sprintf("%s/%s.%s", "/tmp", fn.String(), ext))
}
