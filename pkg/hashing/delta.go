package hashing

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// GenerateDelta ...
func GenerateDelta(origFile io.Reader, sigFile io.Reader, deltaFile io.Writer) error {
	logrus.Debug("delta generation")
	// TODO immplement delta generation
	return fmt.Errorf("not implemented")
}
