package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP HTTP `mapstructure:"http"`
	DB   DB   `mapstructure:"db"`
	JWT  JWT  `mapstructure:"jwt"`
	Log  Log  `mapstructure:"log"`
	CORS CORS `mapstructure:"cors"`
	SMTP SMTP `mapstructure:"smtp"`
}

type HTTP struct {
	Port string `mapstructure:"port"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
	DSN      string `mapstructure:"dsn"`
}

type JWT struct {
	Secret     string        `mapstructure:"secret"`
	AccessTTL  time.Duration `mapstructure:"access_ttl"`
	RefreshTTL time.Duration `mapstructure:"refresh_ttl"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type CORS struct {
	Origins string `mapstructure:"origins"`
}

type SMTP struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("HTTP_PORT", "3000")
	viper.SetDefault("ENV", "local")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("CORS_ORIGINS", "http://localhost:8080,https://turivo-flow.vercel.app")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "turivo")
	viper.SetDefault("DB_PASSWORD", "turivo")
	viper.SetDefault("DB_NAME", "turivo")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "UTC")
	viper.SetDefault("JWT_SECRET", "change-me-in-prod")
	viper.SetDefault("JWT_ACCESS_TTL", "15m")
	viper.SetDefault("JWT_REFRESH_TTL", "168h")
	viper.SetDefault("SMTP_PORT", 465)

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if .env file doesn't exist, we'll use env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Log the error but don't fail - we can use environment variables
			fmt.Printf("Warning: Could not read config file: %v\n", err)
		}
	}

	config := &Config{}

	// Map environment variables to config struct
	config.HTTP.Port = viper.GetString("HTTP_PORT")
	config.DB.Host = viper.GetString("DB_HOST")
	config.DB.Port = viper.GetString("DB_PORT")
	config.DB.User = viper.GetString("DB_USER")
	config.DB.Password = viper.GetString("DB_PASSWORD")
	config.DB.Name = viper.GetString("DB_NAME")
	config.DB.SSLMode = viper.GetString("DB_SSLMODE")
	config.DB.Timezone = viper.GetString("DB_TIMEZONE")
	config.DB.DSN = viper.GetString("DB_DSN")
	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.Log.Level = viper.GetString("LOG_LEVEL")
	config.CORS.Origins = viper.GetString("CORS_ORIGINS")
	config.SMTP.Host = viper.GetString("SMTP_HOST")
	config.SMTP.Port = viper.GetInt("SMTP_PORT")
	config.SMTP.Username = viper.GetString("SMTP_USERNAME")
	config.SMTP.Password = viper.GetString("SMTP_PASSWORD")
	config.SMTP.From = viper.GetString("SMTP_FROM")

	// Parse JWT TTL
	accessTTL, err := time.ParseDuration(viper.GetString("JWT_ACCESS_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}
	config.JWT.AccessTTL = accessTTL

	refreshTTL, err := time.ParseDuration(viper.GetString("JWT_REFRESH_TTL"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TTL: %w", err)
	}
	config.JWT.RefreshTTL = refreshTTL

	// Generate DSN if not provided
	if config.DB.DSN == "" {
		config.DB.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port,
			config.DB.Name, config.DB.SSLMode)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.JWT.Secret == "change-me-in-prod" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}
	if c.DB.DSN == "" {
		return fmt.Errorf("DB_DSN is required")
	}
	if c.SMTP.Host == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.SMTP.Username == "" {
		return fmt.Errorf("SMTP_USERNAME is required")
	}
	if c.SMTP.Password == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}
	return nil
}
