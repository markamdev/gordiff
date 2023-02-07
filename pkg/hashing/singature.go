package hashing

import (
	"encoding/binary"
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
		roller:       NewSimpleRoller(),
	}
}

type sgen struct {
	blockSize    int
	md4hasher    hash.Hash
	weakSumLen   int
	strongSumLen int
	roller       RollSum
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
		s.roller.Init()
		rdCnt, _ := inFile.Read(buffer)
		if rdCnt == 0 {
			// exit due to empty buffer
			break
		}

		// compute strong sum
		mdCnt, err := s.md4hasher.Write(buffer[:rdCnt])
		if mdCnt != rdCnt || err != nil {
			return fmt.Errorf("failed to compute md4: %w", err)
		}

		//compute weak sum
		s.roller.Update(buffer[:rdCnt])
		wSumBuffer := s.roller.Digest()

		// save weak sum for this block to file
		wrCnt, err := sigFile.Write(wSumBuffer)
		if wrCnt != s.weakSumLen || err != nil {
			return fmt.Errorf("failed to write weak sum: %w", err)
		}

		// save strong sum for this block to file
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
	// hardcoded values compatible with RDIFF
	headerBuff[0] = 'r'
	headerBuff[1] = 's'
	headerBuff[2] = 0x01 // signature file
	headerBuff[3] = 0x36 // rollsum + MD4 alghorithms
	// ======
	binary.BigEndian.PutUint32(headerBuff[4:], uint32(s.blockSize))
	binary.BigEndian.PutUint32(headerBuff[8:], uint32(s.strongSumLen))

	n, err := out.Write(headerBuff)
	if n != len(headerBuff) || err != nil {
		return fmt.Errorf("failed to write header magic: %w", err)
	}

	return nil
}
