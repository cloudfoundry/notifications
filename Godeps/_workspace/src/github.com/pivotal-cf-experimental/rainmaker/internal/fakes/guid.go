package fakes

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

func NewGUID(prefix string) string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, guid.String())
	}

	return guid.String()
}
