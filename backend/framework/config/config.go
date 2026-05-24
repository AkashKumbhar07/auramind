package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"database"`
	GRPC   GRPCConfig   `mapstructure:"grpc"`
	Msg    MsgConfig    `mapstructure:"messaging"`
	AI     AIConfig     `mapstructure:"ai"`
	Market MarketConfig `mapstructure:"market"`
	Log    LogConfig    `mapstructure:"logging"`
	Mon    MonConfig    `mapstructure:"monitoring"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

type AppConfig struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
	Port  int    `mapstructure:"port"`
}

type DBConfig struct {
	Postgres  PostgresConfig  `mapstructure:"postgres"`
	ClickHouse ClickHouseConfig `mapstructure:"clickhouse"`
	Redis     RedisConfig     `mapstructure:"redis"`
	SQLite    SQLiteConfig    `mapstructure:"sqlite"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode,
	)
}

type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

type GRPCConfig struct {
	Port          int `mapstructure:"port"`
	MaxRecvMsgSz  int `mapstructure:"max_recv_msg_size"`
	MaxSendMsgSz  int `mapstructure:"max_send_msg_size"`
}

type MsgConfig struct {
	Provider string     `mapstructure:"provider"`
	NATS     NATSConfig `mapstructure:"nats"`
}

type NATSConfig struct {
	URL string `mapstructure:"url"`
}

type AIConfig struct {
	Provider       string  `mapstructure:"provider"`
	Model          string  `mapstructure:"model"`
	EmbeddingModel string  `mapstructure:"embedding_model"`
	Temperature    float64 `mapstructure:"temperature"`
	MaxTokens      int     `mapstructure:"max_tokens"`
}

type MarketConfig struct {
	Exchanges []ExchangeConfig `mapstructure:"exchanges"`
}

type ExchangeConfig struct {
	Name    string `mapstructure:"name"`
	WsURL   string `mapstructure:"ws_url"`
	RestURL string `mapstructure:"rest_url"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MonConfig struct {
	Prometheus PromConfig `mapstructure:"prometheus"`
}

type PromConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Port    int  `mapstructure:"port"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	JWTExpiry string `mapstructure:"jwt_expiry"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetEnvPrefix("AURA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	env := "development"
	path := fmt.Sprintf("configs/environments/%s.yaml", env)
	return Load(path)
}
