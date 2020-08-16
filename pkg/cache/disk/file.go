package disk

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func validateFilename(filename string) error {
	// TODO: need implementation
	return nil
}

func makeTempFile(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	// using the same directory to avoid fs issues like invalid cross-device link
	file, err := ioutil.TempFile(dir, filepath.Base(filename)+".download_*")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func NewTeeFile(r io.ReadCloser, filename string) (io.ReadCloser, error) {
	err := validateFilename(filename)
	if err != nil {
		return nil, err
	}
	tmpFile, err := makeTempFile(filename)
	if err != nil {
		return nil, err
	}
	return &teeFile{
		r:       r,
		tmpFile: tmpFile,
		dst:     filename,
	}, nil
}

type teeFile struct {
	r       io.ReadCloser
	tmpFile *os.File
	dst     string
}

func (t *teeFile) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.tmpFile.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeFile) Close() error {
	err := t.r.Close()
	if err != nil {
		return err
	}
	logrus.Debugln("close reader")
	err = t.tmpFile.Close()
	logrus.Debugln("close tmpFile:", t.tmpFile.Name(), err)
	if err != nil {
		return err
	}
	err = os.Rename(t.tmpFile.Name(), t.dst)
	logrus.Debugln("rename tmpFile:", t.tmpFile.Name(), t.dst, err)
	if err != nil {
		return err
	}
	return nil
}
