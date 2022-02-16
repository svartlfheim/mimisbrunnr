package config

type App interface {
	GetHTTPPort() string
	GetListenHost() string
}

type AppConfig struct {
	HTTPPort string
	ListenHost string
}

func (c *AppConfig) GetHTTPPort() string {
	return c.HTTPPort
}

func (c *AppConfig) GetListenHost() string {
	return c.ListenHost
}

func New() *AppConfig {
	return &AppConfig{
		HTTPPort: "8080",
		ListenHost: "0.0.0.0",
	}
}