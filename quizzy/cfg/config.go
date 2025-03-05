package cfg

import (
	"os"
	"strings"
)

const (
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
	EnvTest        = "TEST"
)

type AppConfig struct {
	// Application environment.
	Env string
	// Address to listen on
	Addr string
	// Firebase configuration file.
	FirebaseConfFile string
	// Base url path.
	BasePath string
	// URI redis
	RedisUri string
}

// getEnvDefault fetch environment variable from the given key and return it if found,
// otherwise the default value is returned.
func getEnvDefault(key, def string) string {
	if v, f := os.LookupEnv(key); f {
		return v
	}

	return def
}

func LoadCfgFromEnv() AppConfig {
	env := strings.ToUpper(getEnvDefault("APP_ENV", EnvProduction))
	
	switch env {
	case EnvDevelopment:
	case EnvProduction:
	case EnvTest:
		break
	default:
		env = EnvProduction
		break
	}

	return AppConfig{
		Env:              env,
		Addr:             getEnvDefault("APP_ADDR", ":8000"),
		FirebaseConfFile: os.Getenv("APP_FIREBASE_CONF_FILE"),
		BasePath:         getEnvDefault("APP_BASE_PATH", "/"),
		RedisUri:         os.Getenv("APP_REDIS_URI"),
	}
}
