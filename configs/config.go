package configs

import (
    "fmt"

    "github.com/spf13/viper"
)

type ServerConfig struct {
    RestPort int `mapstructure:"rest_port"`
    GrpcPort int `mapstructure:"grpc_port"`
}

type StorageConfig struct {
    PostgresDSN string `mapstructure:"postgres_dsn"`
    ClickHouseDSN string `mapstructure:"clickhouse_dsn"`
    RedisAddr string `mapstructure:"redis_addr"`
    MinioEndpoint string `mapstructure:"minio_endpoint"`
    MinioAccessKey string `mapstructure:"minio_access_key"`
    MinioSecretKey string `mapstructure:"minio_secret_key"`
    MinioBucket string `mapstructure:"minio_bucket"`
}

type EventConfig struct {
    KafkaBrokers []string `mapstructure:"kafka_brokers"`
    TopicPrefix string `mapstructure:"topic_prefix"`
}

type LLMConfig struct {
    YandexAPIURL string `mapstructure:"yandex_api_url"`
    YandexAPIKey string `mapstructure:"yandex_api_key"`
}

type TelemetryConfig struct {
    ServiceName string `mapstructure:"service_name"`
    OTelEndpoint string `mapstructure:"otel_endpoint"`
}

type Config struct {
    Server ServerConfig `mapstructure:"server"`
    Storage StorageConfig `mapstructure:"storage"`
    Event EventConfig `mapstructure:"event"`
    LLM LLMConfig `mapstructure:"llm"`
    Telemetry TelemetryConfig `mapstructure:"telemetry"`
}

func Load() (*Config, error) {
    v := viper.New()
    v.SetConfigFile("configs/config.yaml")
    v.AutomaticEnv()
    v.SetEnvPrefix("REPDOC")
    if err := v.ReadInConfig(); err != nil {
        return nil, err
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
