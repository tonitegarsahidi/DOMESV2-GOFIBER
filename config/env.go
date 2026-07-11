package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	DB      DatabaseConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Captcha CaptchaConfig
}

type ServerConfig struct {
	Port              string
	Env               string
	AllowedOrigins    string
	StatsSyncInterval time.Duration
}

type DatabaseConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Charset   string
	ParseTime bool
	Loc       string
}

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type CaptchaConfig struct {
	SecretKey string
	SiteKey   string
	Enabled   bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	redisEnabled, _ := strconv.ParseBool(getEnv("REDIS_ENABLED", "false"))
	recaptchaEnabled, _ := strconv.ParseBool(getEnv("RECAPTCHA_ENABLED", "true"))
	dbParseTime, _ := strconv.ParseBool(getEnv("DB_PARSE_TIME", "true"))

	statsSyncInterval, err := time.ParseDuration(getEnv("STATS_SYNC_INTERVAL", "1h"))
	if err != nil {
		statsSyncInterval = 1 * time.Hour
	}

	return &Config{
		Server: ServerConfig{
			Port:              getEnv("PORT", "3000"),
			Env:               getEnv("ENV", "development"),
			AllowedOrigins:    getEnv("CORS_ALLOWED_ORIGINS", ""),
			StatsSyncInterval: statsSyncInterval,
		},
		DB: DatabaseConfig{
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnv("DB_PORT", "3306"),
			User:      getEnv("DB_USER", "root"),
			Password:  getEnv("DB_PASSWORD", ""),
			Name:      getEnv("DB_NAME", "domesv2"),
			Charset:   getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime: dbParseTime,
			Loc:       getEnv("DB_LOC", "Local"),
		},
		Redis: RedisConfig{
			Enabled:  redisEnabled,
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			ExpiresIn: jwtExpiresIn,
		},
		Captcha: CaptchaConfig{
			SecretKey: getEnv("RECAPTCHA_SECRET_KEY", ""),
			SiteKey:   getEnv("RECAPTCHA_SITE_KEY", ""),
			Enabled:   recaptchaEnabled,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var AppConfig *Config

func InitConfig() {
	AppConfig = LoadConfig()
}
