package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     []byte
	JwtSigninMethod     *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
	CorsConfig
}

func (c *Config) readConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found, using default values")
	}

	maxAge, err := strconv.Atoi(getEnv("CORS_MAX_AGE", "3600"))
	if err != nil {
		maxAge = 3600
	}

	accessTokenLifeTime, err := strconv.Atoi(getEnv("ACCESS_TOKEN_LIFETIME", "3600"))
	if err != nil {
		accessTokenLifeTime = 3600
	}

	c.DBConfig = DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		Database: getEnv("DB_NAME", "go_api_library"),
		Username: getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Driver:   getEnv("DB_DRIVER", "mysql"),
	}

	c.APIConfig = APIConfig{
		ApiPort: getEnv("API_PORT", "8888"),
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:     getEnv("APP_NAME", "GO API LIBRARY"),
		JwtSignatureKey:     []byte(getEnv("JWT_SECRET", "random-secret-key")),
		JwtSigninMethod:     jwt.SigningMethodHS256,
		AccessTokenLifeTime: time.Duration(accessTokenLifeTime) * time.Second,
	}

	c.CorsConfig = CorsConfig{
		AllowOrigins:     []string{getEnv("CORS_ALLOW_ORIGINS", "*")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: getEnv("CORS_ALLOW_CREDENTIALS", "true") == "true",
		MaxAge:           maxAge,
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
