package config

import "os"

type Config struct {
	ServerHost      string
	PostgresConnect string
	JwtSecretKey    string
	LlmRequestLink  string
	LlmHeaders      string
}

func GetConfig() Config {
	return Config{
		ServerHost:      os.Getenv("SERVER_HOST"),
		PostgresConnect: os.Getenv("POSTGRES_CONNECT"),
		JwtSecretKey:    os.Getenv("JWT_SECRET_KEY"),
		LlmRequestLink:  os.Getenv("LM_REQUEST_LINK"),
		LlmHeaders:      os.Getenv("LM_HEADERS"),
	}
}
