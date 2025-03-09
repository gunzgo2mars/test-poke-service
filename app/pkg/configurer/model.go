package configurer

import (
	"time"
)

type AppConfig struct {
	App     App   `mapstructure:"app"`
	Log     Log   `mapstructure:"log"`
	MySQL   MySQL `mapstructure:"mysql"`
	Redis   Redis `mapstructure:"redis"`
	Http    Http  `mapstructure:"http"`
	Secrets Secrets
}

type App struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type Log struct {
	Env string `mapstructure:"env"`
}

type MySQL struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type Redis struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type Http struct {
	PokeAPI HttpPokeAPI `mapstructure:"pokeapi"`
}

type HttpPokeAPI struct {
	BaseUrl          string        `mapstructure:"base_url"`
	RetryWaitTime    time.Duration `mapstructure:"retry_wait_time"`
	RetryMaxWaitTime time.Duration `mapstructure:"retry_max_wait_time"`
	RetryCount       int           `mapstructure:"retry_count"`
	Timeout          time.Duration `mapstructure:"timeout"`
}

type Secrets struct {
	MySQLUser     string `envconfig:"SECRET_MYSQL_USER"`
	MySQLPassword string `envconfig:"SECRET_MYSQL_PASSWORD"`
	MySQLDBName   string `envconfig:"SECRET_MYSQL_DBNAME"`
	JwtKey        string `envconfig:"SECRET_JWT_KEY"`
}
