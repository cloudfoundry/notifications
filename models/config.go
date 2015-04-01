package models

type Config struct {
	DatabaseURL         string
	MigrationsPath      string
	DefaultTemplatePath string
	MaxOpenConnections  int
}
