package config

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type Config struct {
	Verbose bool
	Addr    string
	Rewrite *Rewrite
	Cache   *Cache
}

type Rewrite struct {
}

type Cache struct {
	Dir   string
	Rules []*Rule
}

type Rule struct {
	Name       string
	Conditions []map[string]string
}

var defaultConfig = &Config{
	Cache: &Cache{
		Dir:   "",
		Rules: nil,
	},
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

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	// TODO: validate config
	cfg := defaultConfig
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalln("err when load config:", cfg)
	}
	if cfg.Verbose {
		spew.Dump(cfg)
	}
	return cfg
}
