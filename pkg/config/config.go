package config

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var (
	validate      = validator.New()
	defaultConfig = &Config{
		Rewrite: &Rewrite{
			Rules: []*RewriteRule{},
		},
		Cache: &Cache{
			Dir:   "",
			Rules: nil,
		},
	}
)

type Config struct {
	Verbose bool
	Addr    string
	Rewrite *Rewrite
	Cache   *Cache
}

func LoadConfig(cfgFile string) *Config {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".mirroroman")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("err when read config file:", err)
	}
	cfg := defaultConfig
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalln("err when unmarshal config:", cfgFile)
	}
	if cfg.Verbose {
		spew.Dump(cfg)
	}
	err = validate.Struct(cfg)
	if err != nil {
		log.Fatalln("invalid config:", err)
	}
	return cfg
}
