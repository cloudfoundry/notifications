package services

type TokenLoader interface {
	Load() (string, error)
}
