package hashing

import (
	"io"

	"github.com/sirupsen/logrus"
)

// GenerateDelta ...
func GenerateDelta(origFile io.Reader, sigFile io.Reader, deltaFile io.Writer) error {
	logrus.Debug("delta generation")
	// TODO immplement delta generation
	buffer := make([]byte, 0, 1024)
	var err error
	var cnt int
	for err == nil {
		cnt, err = origFile.Read(buffer)
		logrus.Debugf("cnt: %d err: %v", cnt, err)
		if cnt > 0 {
			deltaFile.Write(buffer[:cnt])
		}
	}

	return nil
}
