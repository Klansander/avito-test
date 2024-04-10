//go:build !windows

package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App  struct{} `yaml:"app"`
	HTTP struct {
		Host              string        `yaml:"host" env:"HTTP_HOST" env-required:"true"`
		Port              int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
		ReadTimeout       time.Duration `yaml:"read-timeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
		WriteTimeout      time.Duration `yaml:"write-timeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
		IdleTimeout       time.Duration `yaml:"idle-timeout" env:"HTTP_IDLE_TIMEOUT" env-required:"true"`
		ReadHeaderTimeout time.Duration `yaml:"read-header-timeout" env:"HTTP_READ_HEADER_TIMEOUT" env-required:"true"`
		ReRequest         time.Duration `yaml:"re-request" env:"HTTP_RE_REQUEST" env-required:"true"`
		CORS              struct {
			AllowedMethods     []string `yaml:"allowed-methods" env:"CORS_ALLOWED_METHODS" env-required:"true"`
			AllowedOrigins     []string `yaml:"allowed-origins" env:"CORS_ALLOWED_ORIGINS" env-required:"true"`
			AllowCredentials   bool     `yaml:"allow-credentials" env:"CORS_ALLOW_CREDENTIALS" env-required:"true"`
			AllowedHeaders     []string `yaml:"allowed-headers" env:"CORS_ALLOWED_HEADERS" env-required:"true"`
			OptionsPassthrough bool     `yaml:"options-passthrough" env:"CORS_OPTIONS_PASSTHROUGH" env-required:"true"`
			ExposedHeaders     []string `yaml:"exposed-headers" env:"CORS_EXPOSED_HEADERS" env-required:"false"`
			Debug              bool     `yaml:"debug" env:"CORS_DEBUG" env-required:"true"`
		} `yaml:"cors"`
		StartFront bool   `yaml:"start-front" env:"HTTP_START_FRONT"`
		DistFolder string `yaml:"dist-folder" env:"HTTP_DIST_FOLDER"`
		DistPort   int    `yaml:"dist-port" env:"HTTP_DIST_PORT"`
	} `yaml:"http"`
	Cron struct {
		Interval time.Duration `yaml:"CRON_INTERVAL"`
	} `yaml:"cron"`
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
		Port int    `yaml:"port" env:"REDIS_PORT" env-required:"true"`
	} `yaml:"redis"`

	PSQL struct {
		Username string        `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
		Password string        `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
		Host     string        `yaml:"host" env:"PSQL_HOST" env-required:"true"`
		Port     int           `yaml:"port" env:"PSQL_PORT" env-required:"true"`
		Database string        `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
		Timeout  time.Duration `yaml:"timeout" env:"PSQL_TIMEOUT" env-required:"true"`
		LimitMax int           `yaml:"limit-max" env:"PSQL_LIMIT_MAX" env-required:"true"`
	} `yaml:"psql"`
}

const configPath = "config.local.yaml"

var once sync.Once
var cfg *Config

func New() *Config {

	once.Do(func() {

		cfg = &Config{}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			if err := cleanenv.ReadEnv(cfg); err != nil {
				helpText := "Список переменных окружения"
				help, _ := cleanenv.GetDescription(cfg, &helpText)
				logrus.Print(help)
				logrus.Fatal(err)
			}
		}

	})

	return cfg

}
