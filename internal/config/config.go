package config

import (
	"fmt"
	"github.com/spf13/viper"
)


type Config struct {
	HTTPPort string
	Postgres PostgresConfig
}


type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}


func (p *PostgresConfig) DSN() string {
	// host=db port=5432 user=user password=password dbname=subscriptions_db sslmode=disable
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}


func LoadConfig(path string) (*Config, error) {

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")


	viper.AutomaticEnv()


	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
		}
	}

	
	cfg := &Config{
		HTTPPort: viper.GetString("HTTP_PORT"),
		Postgres: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
		},
	}
	

	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}

	return cfg, nil
}
