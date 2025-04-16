package config

import (
	"time"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Config struct {
	ServerPort    string
	Database      DatabaseConfig
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "12345678")
	viper.SetDefault("DB_NAME", "authforge")

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("REFRESH_EXPIRY", "168h")

	if err := viper.ReadInConfig(); err != nil {
	}

	cfg := &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		SMTPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:      viper.GetInt("SMTP_PORT"),
		SMTPUsername:  viper.GetString("SMTP_USERNAME"),
		SMTPPassword:  viper.GetString("SMTP_PASSWORD"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		JWTExpiry:     viper.GetDuration("JWT_EXPIRY"),
		RefreshExpiry: viper.GetDuration("REFRESH_EXPIRY"),
	}
	return cfg, nil
}
