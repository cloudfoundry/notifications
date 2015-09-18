package models

import (
	"fmt"
	"io"
)

type GUIDGenerator struct {
	reader io.Reader
}

func NewGUIDGenerator(reader io.Reader) GUIDGenerator {
	return GUIDGenerator{
		reader: reader,
	}
}

func (g GUIDGenerator) Generate() (string, error) {
	var buf [16]byte

	_, err := g.reader.Read(buf[:])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:]), nil
}
