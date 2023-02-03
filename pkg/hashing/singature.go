package hashing

import (
	"fmt"
	"hash"
	"io"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/md4"
)

// SigGen ...
type SigGen interface {
	Compute(io.Reader, io.Writer) error
}

// NewSigGen ...
func NewSigGen(blockSize int) SigGen {
	return &sgen{
		blockSize:    blockSize,
		md4hasher:    md4.New(),
		weakSumLen:   4,
		strongSumLen: 8,
	}
}

type sgen struct {
	blockSize    int
	md4hasher    hash.Hash
	weakSumLen   int
	strongSumLen int
}

// GenerateSignature ...
func (s *sgen) Compute(inFile io.Reader, sigFile io.Writer) error {
	logrus.Debug("signature generation")

	err := s.writeHeader(sigFile)
	if err != nil {
		return fmt.Errorf("signature computation failed: %w", err)
	}

	buffer := make([]byte, s.blockSize)
	for {
		logrus.Debug("- new block processing")
		s.md4hasher.Reset()
		rdCnt, _ := inFile.Read(buffer)
		if rdCnt == 0 {
			// exit due to empty buffer
			break
		}

		mdCnt, err := s.md4hasher.Write(buffer[:rdCnt])
		if mdCnt != rdCnt || err != nil {
			return fmt.Errorf("failed to compute md4: %w", err)
		}

		// save weak sum
		wSumBuffer := make([]byte, s.weakSumLen)
		wrCnt, err := sigFile.Write(wSumBuffer)
		if wrCnt != s.weakSumLen || err != nil {
			return fmt.Errorf("failed to write weak sum: %w", err)
		}

		// save storng sum
		mdSum := s.md4hasher.Sum(nil)
		wrCnt, err = sigFile.Write(mdSum[:s.strongSumLen]) // need only first half of computed sum

		if wrCnt != s.strongSumLen || err != nil {
			return fmt.Errorf("failed to write strong sum: %w", err)
		}
		if rdCnt != s.blockSize {
			// exit as not-full buffer means end of data
			break
		}
	}

	return nil
}

func (s *sgen) writeHeader(out io.Writer) error {
	logrus.Debug("- header writing")

	headerBuff := make([]byte, 12)
	headerBuff[0] = 'M'
	headerBuff[1] = 'K'

	n, err := out.Write(headerBuff)
	if n != len(headerBuff) || err != nil {
		return fmt.Errorf("failed to write header magic: %w", err)
	}

	return nil
}

func generateMD4() string {
	md4obj := md4.New()

	return fmt.Sprint("BlockSize:", md4obj.BlockSize(), "SumSize:", md4obj.Size())
}
