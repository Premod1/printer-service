package config

type Config struct {
	WebSocketPort string
	HTTPPort      string
}

func Load() *Config {
	return &Config{
		WebSocketPort: ":8081",
		HTTPPort:      ":8080",
	}
}
