package toolbelt

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/mkorenkov/sparsed"
)

//IsGzip checks whether the given stream is gzipped.
func IsGzip(r io.ReadSeeker) (bool, error) {
	bufferedReader := bufio.NewReader(r)
	testBytes, err := bufferedReader.Peek(2) //read 2 bytes
	if err != nil {
		return false, err
	}
	// disregarding of reader.Peek docs, need to ofset the file back to beginning..
	_, err = r.Seek(0, 0)
	if err != nil {
		return false, err
	}
	// http://stackoverflow.com/questions/28309988/how-to-read-from-either-gzip-or-plain-text-reader-in-golang
	return testBytes[0] == 31 && testBytes[1] == 139, nil
}

//UntarFilename extracts archive at given path.
func UntarFilename(filename string, destinationDir string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return UntarFile(file, destinationDir)
}

//UntarFile extracts archived file.
func UntarFile(file *os.File, destinationDir string) error {
	isGzipped, err := IsGzip(file)
	if err != nil {
		return err
	}

	fileReader := io.ReadCloser(file)
	if isGzipped {
		fileReader, err = gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer fileReader.Close()
	}

	err = UntarStream(fileReader, destinationDir)
	if err != nil {
		return err
	}
	return nil
}

//UntarStream extracts tar from io.Reader.
func UntarStream(fileReader io.Reader, destinationDir string) error {
	tarBallReader := tar.NewReader(fileReader)
	for {
		header, err := tarBallReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		filename := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			// directory
			err = os.MkdirAll(fmt.Sprintf("%v/%v", destinationDir, filename), os.FileMode(header.Mode))
			if err != nil {
				return err
			}

		case tar.TypeReg:
			// file
			err := func() error {
				w, err := os.OpenFile(fmt.Sprintf("%v/%v", destinationDir, filename), os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
				if err != nil {
					return err
				}
				writer := sparsed.NewSparseFilesWriter(w)
				defer func() {
					if writer != nil {
						writer.Flush()
					}
					w.Close()
				}()

				io.Copy(writer, tarBallReader)
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
