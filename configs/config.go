package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	RestPort int `mapstructure:"rest_port"`
	GrpcPort int `mapstructure:"grpc_port"`
}

type StorageConfig struct {
	PostgresDSN    string `mapstructure:"postgres_dsn"`
	ClickHouseDSN  string `mapstructure:"clickhouse_dsn"`
	RedisAddr      string `mapstructure:"redis_addr"`
	MinioEndpoint  string `mapstructure:"minio_endpoint"`
	MinioAccessKey string `mapstructure:"minio_access_key"`
	MinioSecretKey string `mapstructure:"minio_secret_key"`
	MinioBucket    string `mapstructure:"minio_bucket"`
}

type EventConfig struct {
	KafkaBrokers []string `mapstructure:"kafka_brokers"`
	TopicPrefix  string   `mapstructure:"topic_prefix"`
}

type LLMConfig struct {
	YandexAPIURL string `mapstructure:"yandex_api_url"`
	YandexAPIKey string `mapstructure:"yandex_api_key"`
}

type TelemetryConfig struct {
	ServiceName  string `mapstructure:"service_name"`
	OTelEndpoint string `mapstructure:"otel_endpoint"`
}

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Event     EventConfig     `mapstructure:"event"`
	LLM       LLMConfig       `mapstructure:"llm"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("REPDOC")

	configCandidates := []string{
		os.Getenv("REPDOC_CONFIG_FILE"),
		"/configs/config.yaml",
		"configs/config.yaml",
	}
	var readErr error
	for _, path := range configCandidates {
		if path == "" {
			continue
		}
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err == nil {
			readErr = nil
			break
		} else {
			readErr = err
		}
	}
	if readErr != nil {
		return nil, readErr
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) RestAddress() string {
	return fmt.Sprintf(":%d", c.Server.RestPort)
}

func (c Config) GrpcAddress() string {
	return fmt.Sprintf(":%d", c.Server.GrpcPort)
}
