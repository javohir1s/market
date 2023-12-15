package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	Error = "error >>> "
	Info  = "info >>> "
	Log   = "log >>> "
)

type Config struct {
	PostgresHost          string
	PostgresUser          string
	PostgresDatabase      string
	PostgresPassword      string
	PostgresPort          string
	PostgresMaxConnection int32

	RedisHost     string
	RedisPort     string
	RedisDatabase int32
	RedisPassword string

	ServiceHost     string
	ServiceHTTPPort string

	SecretKey string
}

func Load() Config {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("not found env")
	}

	var cfg Config

	cfg.ServiceHost = cast.ToString(getValueOrDefault("SERVICE_HOST", "localhost"))
	cfg.ServiceHTTPPort = cast.ToString(getValueOrDefault("SERVICE_HTTP_PORT", ":8080"))

	cfg.PostgresHost = cast.ToString(getValueOrDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresUser = cast.ToString(getValueOrDefault("POSTGRES_USER", "javohir"))
	cfg.PostgresDatabase = cast.ToString(getValueOrDefault("POSTGRES_DATABASE", "market"))
	cfg.PostgresPassword = cast.ToString(getValueOrDefault("POSTGRES_PASSWORD", "12345"))
	cfg.PostgresPort = cast.ToString(getValueOrDefault("POSTGRES_PORT", "5432"))
	cfg.PostgresMaxConnection = cast.ToInt32(getValueOrDefault("POSTGRES_MAX_CONN", 30))

	cfg.RedisHost = cast.ToString(getValueOrDefault("REDIS_HOST", "localhost"))
	cfg.RedisPort = cast.ToString(getValueOrDefault("REDIS_PORT", ":6379"))
	cfg.RedisDatabase = cast.ToInt32(getValueOrDefault("REDIS_DB", 0))
	cfg.RedisPassword = cast.ToString(getValueOrDefault("REDIS_PASSWORD", ""))

	cfg.SecretKey = cast.ToString(getValueOrDefault("SECRET_KEY", "q6T6LlwdRk"))

	return cfg
}

func getValueOrDefault(key string, defaultValue interface{}) interface{} {

	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
