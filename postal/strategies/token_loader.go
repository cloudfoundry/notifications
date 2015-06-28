package strategies

type TokenLoader interface {
	Load() (string, error)
}
