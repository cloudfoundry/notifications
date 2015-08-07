package fakes

import (
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/nu7hatch/gouuid"
)

var GUIDGenerator = services.GUIDGenerationFunc(func() (*uuid.UUID, error) {
	guid := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
	return &guid, nil
})

type IncrementingGUIDGenerator struct {
	increment byte
}

func NewIncrementingGUIDGenerator() *IncrementingGUIDGenerator {
	return &IncrementingGUIDGenerator{
		increment: byte(0),
	}
}

func (i *IncrementingGUIDGenerator) Generate() (*uuid.UUID, error) {
	guid := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55 + i.increment})
	i.increment = i.increment + byte(1)
	return &guid, nil
}
