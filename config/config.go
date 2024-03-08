package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresDatabase string
	PostgresPassword string

	RedisHost string
	RedisPort int

	UserServiceHost string
	UserServicePort int

	HealthcareServiceHost string
	HealthcareServicePort int

	BookingServiceHost string
	BookingServicePort int

	CtxTimeout int

	LogLevel string
	HTTPPort string

	SignInKey           string
	AccessTokenTimeOut  int
	RefreshTokenTimeOut int

	AuthConfigPath string

	SendEmailFrom string
	EmailCode     string

	AccessTokenTimeout  int //minutes
	RefreshTokenTimeout int //hours
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "auth"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "mubina2007"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "localhost"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 8080))

	c.HealthcareServiceHost = cast.ToString(getOrReturnDefault("HEALTHCARE_SERVICE_HOST", "localhost"))
	c.HealthcareServicePort = cast.ToInt(getOrReturnDefault("HEALTHCARE_SERVICE_PORT", 6060))

	c.BookingServiceHost = cast.ToString(getOrReturnDefault("BOOKING_SERVICE_HOST", "localhost"))
	c.BookingServicePort = cast.ToInt(getOrReturnDefault("BOOKING_SERVICE_PORT", 9091))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	c.AccessTokenTimeout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 500))
	c.RefreshTokenTimeout = cast.ToInt(getOrReturnDefault("REFRESH_TOKEN_TIMEOUT", 3))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":7070"))

	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "abc"))
	c.AccessTokenTimeOut = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 2000))
	c.RefreshTokenTimeOut = cast.ToInt(getOrReturnDefault("REFRESH_TOKEN_TIMEOUT", 3000))

	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))

	c.SendEmailFrom = cast.ToString(getOrReturnDefault("EMAIL_FROM", "mubinayigitaliyeva00@gmail.com"))
	c.EmailCode = cast.ToString(getOrReturnDefault("EMAIL_CODE", "iocd vnhb lnvx digm"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
