package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	WebPush  WebPushConfig
	DataRelay DataRelayConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	Port string
	Host string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type WebPushConfig struct {
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubject    string
}

type DataRelayConfig struct {
	URL   string
	Token string
}

type RabbitMQConfig struct {
	URL                string
	QueueNotifications string
	Workers            int
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Host: viper.GetString("SERVER_HOST"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		WebPush: WebPushConfig{
			VAPIDPublicKey:  viper.GetString("VAPID_PUBLIC_KEY"),
			VAPIDPrivateKey: viper.GetString("VAPID_PRIVATE_KEY"),
			VAPIDSubject:    viper.GetString("VAPID_SUBJECT"),
		},
		DataRelay: DataRelayConfig{
			URL:   viper.GetString("DATA_RELAY_API_URL"),
			Token: viper.GetString("DATA_RELAY_API_TOKEN"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:                viper.GetString("RABBITMQ_URL"),
			QueueNotifications: viper.GetString("RABBITMQ_QUEUE_NOTIFICATIONS"),
			Workers:            viper.GetInt("RABBITMQ_WORKERS"),
		},
	}

	return config, nil
}
