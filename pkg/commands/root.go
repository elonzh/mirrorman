package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/elonzh/mirrorman/pkg/config"
)

type Executor struct {
	cfg     *config.Config
	rootCmd *cobra.Command
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

func NewExecutor() *Executor {
	var cfgFile string
	e := &Executor{
		rootCmd: &cobra.Command{
			Use: "mirroroman",
		},
	}
	e.rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mirroroman.yaml)")
	e.rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	err := viper.BindPFlag("verbose", e.rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		panic(err)
	}
	e.initServeCommand()
	cobra.OnInitialize(func() {
		cfgFile, err := e.rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			log.Fatalln("err when get flag `config`:", e)
		}
		e.cfg = config.LoadConfig(cfgFile)
	})
	return e
}
