package hashing

import (
	"io"

	"github.com/sirupsen/logrus"
)

// GenerateSignature ...
func GenerateSignature(inFile io.Reader, sigFile io.Writer) error {
	logrus.Debug("signature generation")
	// TODO immplement signature generation
	buffer := make([]byte, 0, 1024)
	sigFile.Write([]byte("signature:"))
	for {
		cnt, err := inFile.Read(buffer)
		logrus.Debugf("cnt: %d err: %v", cnt, err)
		if cnt == 0 {
			break
		}
		sigFile.Write(buffer[:cnt])
	}

	return nil
}
