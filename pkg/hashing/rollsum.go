package hashing

/*
This file is a simplified port of RollingSum from librsync used in rdiff tool.

Original implementation is available in librsync repository:
https://github.com/librsync/librsync/

Porting made for compatibility with rdiff output
*/

const (
	rollsumOffset = 31
)

// RollSum ...
type RollSum interface {
	Init()
	Update([]byte)
	Digest() []byte
}

// NewSimpleRoller ...
func NewSimpleRoller() RollSum {
	return &simpleRoller{count: 0, s1: 0, s2: 0}
}

type simpleRoller struct {
	count int
	s1    uint16
	s2    uint16
}

func (sr *simpleRoller) Init() {
	sr.count = 0
	sr.s1 = 0
	sr.s2 = 0
}

func (sr *simpleRoller) Update(inBuffer []byte) {
	inLen := len(inBuffer)
	tempS1 := uint64(sr.s1)
	tempS2 := uint64(sr.s2)

	for _, input := range inBuffer {
		tempS1 += uint64(input)
		tempS2 += tempS1
	}

	sr.count += inLen
	tempS1 += uint64(inLen) * rollsumOffset
	tempS2 += uint64((inLen * (inLen + 1) / 2) * rollsumOffset)
	sr.s1 = uint16(tempS1)
	sr.s2 = uint16(tempS2)
}

func (sr *simpleRoller) Digest() []byte {
	result := make([]byte, 4)
	result[0] = uint8((sr.s2 >> 8) & 0x0ff)
	result[1] = uint8(sr.s2 & 0x00ff)
	result[2] = uint8((sr.s1 >> 8) & 0x0ff)
	result[3] = uint8(sr.s1 & 0x00ff)

	return result
}
