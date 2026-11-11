package config

import "os"

type (
	Config struct {
		Storage    db
		Connectors connectors
		Addr       addresses
	}

	db struct {
		User string
		Pass string
		DB   string
		Host string
		Port string
	}

	connectors struct {
		Telegram string
	}

	addresses struct {
		HttpAddr     string
		CoreUrl      string
		CoreToken    string
		ContextToken string
	}
)

func getDBConfig() db {
	return db{
		User: os.Getenv("POSTGRES_USER"),
		Pass: os.Getenv("POSTGRES_PASSWORD"),
		DB:   os.Getenv("POSTGRES_DB"),
		Host: os.Getenv("POSTGRES_HOST"),
		Port: os.Getenv("POSTGRES_PORT"),
	}
}

func getConnectorsConfig() connectors {
	return connectors{
		Telegram: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func getAddresses() addresses {
	return addresses{
		HttpAddr:     os.Getenv("HTTP_ADDR"),
		CoreUrl:      os.Getenv("CORE_API_URL"),
		CoreToken:    os.Getenv("CORE_API_TOKEN"),
		ContextToken: os.Getenv("CORE_CONTEXT_TOKEN"),
	}
}

func MustLoad() Config {
	return Config{
		Storage:    getDBConfig(),
		Connectors: getConnectorsConfig(),
		Addr:       getAddresses(),
	}
}
