package mocks

import "gopkg.in/gorp.v1"

type GobbleInitializer struct {
	InitializeDBMapCall struct {
		Receives struct {
			DbMap *gorp.DbMap
		}
	}
}

func NewGobbleInitializer() *GobbleInitializer {
	return &GobbleInitializer{}
}

func (m *GobbleInitializer) InitializeDBMap(dbmap *gorp.DbMap) {
	m.InitializeDBMapCall.Receives.DbMap = dbmap
}
