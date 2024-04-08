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
			AllowedMethods     []string `yaml:"allowed_methods"`
			AllowedOrigins     []string `yaml:"allowed_origins"`
			AllowCredentials   bool     `yaml:"allow_credentials"`
			AllowedHeaders     []string `yaml:"allowed_headers"`
			OptionsPassthrough bool     `yaml:"options_passthrough"`
			ExposedHeaders     []string `yaml:"exposed_headers"`
			Debug              bool     `yaml:"debug"`
		} `yaml:"cors"`
		StartFront bool   `yaml:"start-front" env:"HTTP_START_FRONT"`
		DistFolder string `yaml:"dist-folder" env:"HTTP_DIST_FOLDER"`
		DistPort   int    `yaml:"dist-port" env:"HTTP_DIST_PORT"`
	} `yaml:"http"`

	PSQL struct {
		Username string        `yaml:"username" `
		Password string        `yaml:"password" `
		Host     string        `yaml:"host" `
		Port     int           `yaml:"port" `
		Database string        `yaml:"database" `
		Timeout  time.Duration `yaml:"timeout" env:"PSQL_TIMEOUT" env-required:"true"`
		LimitMax int           `yaml:"limit-max" env:"PSQL_LIMIT_MAX" env-required:"true"`
	} `yaml:"psql"`
}

const configPath = "./config.local.yaml"

var once sync.Once
var cfg *Config

func New() *Config {

	once.Do(func() {

		cfg = &Config{}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			helpText := "Список переменных окружения"
			help, _ := cleanenv.GetDescription(cfg, &helpText)
			logrus.Print(help)
			logrus.Fatal(err)
		}

	})

	return cfg

}
