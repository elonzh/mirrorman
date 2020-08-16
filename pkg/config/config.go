package config

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	validate      = validator.New()
	defaultConfig = &Config{
		HttpAddr:  ":8080",
		ProxyAddr: ":8081",
		Rewrite: &Rewrite{
			Rules: []*RewriteRule{},
		},
		Cache: &Cache{
			Backend: "disk",
			Rules:   nil,
		},
	}
)

type Config struct {
	Verbose   bool
	HttpAddr  string
	ProxyAddr string
	Rewrite   *Rewrite
	Cache     *Cache
}

func LoadConfig(cfgFile string) *Config {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalln(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".mirroroman")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatalln("error when read config file")
	}
	cfg := defaultConfig
	err := viper.Unmarshal(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("error when unmarshal config")
	}
	if cfg.Verbose {
		spew.Dump(cfg)
	}
	err = validate.Struct(cfg)
	if err != nil {
		logrus.WithError(err).Fatalln("invalid config")
	}
	return cfg
}
