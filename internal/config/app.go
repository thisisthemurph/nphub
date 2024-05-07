package config

type AppConfig struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Name     string
	Path     string
	FullPath string
}

func NewAppConfig(getenv func(string) string) *AppConfig {
	return &AppConfig{
		Database: DatabaseConfig{
			Name:     getenv("DB_NAME"),
			Path:     getenv("DB_PATH"),
			FullPath: getenv("DB_PATH") + getenv("DB_NAME"),
		},
	}
}
