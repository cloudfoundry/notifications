package mocks

import "errors"

type GUIDGenerator struct {
	GenerateCall struct {
		CallCount int
		Returns   struct {
			GUIDs []string
			Error error
		}
	}
}

func NewGUIDGenerator() *GUIDGenerator {
	return &GUIDGenerator{}
}

func (g *GUIDGenerator) Generate() (string, error) {
	if g.GenerateCall.CallCount >= len(g.GenerateCall.Returns.GUIDs) {
		return "", errors.New("no GUIDs to return")
	}

	guid := g.GenerateCall.Returns.GUIDs[g.GenerateCall.CallCount]
	g.GenerateCall.CallCount++

	return guid, g.GenerateCall.Returns.Error
}
