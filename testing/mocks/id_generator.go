package mocks

import "errors"

type IDGenerator struct {
	GenerateCall struct {
		CallCount int
		Returns   struct {
			IDs   []string
			Error error
		}
	}
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) Generate() (string, error) {
	if g.GenerateCall.CallCount >= len(g.GenerateCall.Returns.IDs) {
		return "", errors.New("no IDs to return")
	}

	guid := g.GenerateCall.Returns.IDs[g.GenerateCall.CallCount]
	g.GenerateCall.CallCount++

	return guid, g.GenerateCall.Returns.Error
}
