package files

import (
	"githubcom/markamdev/gordiff/pkg/common"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// GetInputFile ...
func GetInputFile(path string) (*os.File, error) {
	logrus.Debugf("opening input file '%s'", path)
	if path == "--" {
		return os.Stdin, nil
	}

	fl, err := os.Open(path)
	if err != nil {
		pathErr := err.(*os.PathError)
		switch pathErr.Err {
		case os.ErrNotExist:
			logrus.Error("file does not exist")
			return nil, common.ErrInvalidFilePath
		case os.ErrPermission:
			logrus.Error("cannot access file")
			return nil, common.ErrCannotAccess
		default:
			logrus.Errorf("unexpected error '%v' occured", err)
			return nil, common.ErrUnexpected
		}
	}
	fl.Seek(0, 0)
	fl.SetReadDeadline(time.Time{})

	return fl, nil
}

// GetOutputFile ...
func GetOutputFile(path string, overwrite bool) (*os.File, error) {
	logrus.Debugf("opening output file '%s", path)
	if path == "--" {
		return os.Stdout, nil
	}

	_, err := os.Stat(path)
	if err == nil && !overwrite {
		// there should be error returned
		logrus.Error("file already exists") // it can be also dir but this has no meaning
		return nil, common.ErrFileExists
	}
	fl, err := os.Create(path)
	if err != nil {
		logrus.Errorf("failed to create output file: %s", err.Error())
	}

	return fl, nil
}
