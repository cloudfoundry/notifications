package mocks

import (
	"errors"

	"github.com/nu7hatch/gouuid"
)

type GUIDGenerator struct {
	GenerateCall struct {
		CallCount int
		Returns   struct {
			GUIDs []*uuid.UUID
			Error error
		}
	}
}

func NewGUIDGenerator() *GUIDGenerator {
	return &GUIDGenerator{}
}

func (g *GUIDGenerator) Generate() (*uuid.UUID, error) {
	if g.GenerateCall.CallCount >= len(g.GenerateCall.Returns.GUIDs) {
		return nil, errors.New("no GUIDs to return")
	}

	guid := g.GenerateCall.Returns.GUIDs[g.GenerateCall.CallCount]
	g.GenerateCall.CallCount++

	return guid, g.GenerateCall.Returns.Error
}
