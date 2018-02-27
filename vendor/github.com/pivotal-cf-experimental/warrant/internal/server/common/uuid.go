package common

import (
	"crypto/rand"
	"fmt"
)

func NewUUID() (string, error) {
	u := [16]byte{}
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
