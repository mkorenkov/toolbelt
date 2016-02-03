package toolbelt

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

//IsGzip checks whether the given stream is gzipped.
func IsGzip(r io.Reader) (bool, error) {
	// http://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang
	bufferedReader := bufio.NewReader(r)
	testBytes, err := bufferedReader.Peek(2) //read 2 bytes
	if err != nil {
		return false, err
	}
	return testBytes[0] == 31 && testBytes[1] == 139, nil
}

//UntarFilename extracts archive at given path.
func UntarFilename(sourcefile string) error {
	file, err := os.Open(sourcefile)
	if err != nil {
		return err
	}
	defer file.Close()
	return UntarFile(file)
}

//UntarFile extracts archived file.
func UntarFile(file *os.File) error {
	isGzipped, err := IsGzip(file)
	if err != nil {
		return err
	}

	var fileReader io.ReadCloser
	if isGzipped {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()
	} else {
		fileReader = io.ReadCloser(file)
	}

	err = UntarStream(fileReader)
	if err != nil {
		return err
	}
	return nil
}

//UntarStream extracts tar from io.Reader.
func UntarStream(fileReader io.Reader) error {
	tarBallReader := tar.NewReader(fileReader)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// directory
			err = os.MkdirAll(filename, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

		case tar.TypeReg:
			// file
			err := func() error {
				writer, err := os.Create(filename)
				if err != nil {
					return err
				}
				defer writer.Close()

				io.Copy(writer, tarBallReader)
				err = os.Chmod(filename, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				return nil
			}()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%c in file %s is not supported yet", header.Typeflag, filename)
		}
	}
	return nil
}
