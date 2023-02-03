package hashing

import (
	"io"

	"github.com/sirupsen/logrus"
)

// GenerateDelta ...
func GenerateDelta(origFile io.Reader, sigFile io.Reader, deltaFile io.Writer) error {
	logrus.Debug("delta generation")
	// TODO immplement delta generation
	buffer := make([]byte, 1024)
	for {
		cnt, _ := origFile.Read(buffer)
		if cnt == 0 {
			break
		}
		deltaFile.Write(buffer[:cnt])
	}

	return nil
}
