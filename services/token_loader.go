package services

type TokenLoaderInterface interface {
	Load(string) (string, error)
}
