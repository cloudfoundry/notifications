package domain

import (
	"math/rand"
	"time"
)

const (
	origin = "uaa"
	schema = "urn:scim:schemas:core:1.0"
)

var schemas = []string{schema}

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
