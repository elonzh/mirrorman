package commands

import (
	"fmt"
	"path"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"

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
			logrus.WithError(err).Fatalln("error when get flag `config`")
		}
		e.cfg = config.LoadConfig(cfgFile)
		logrus.SetReportCaller(true)
		logrus.SetFormatter(&zt_formatter.ZtFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			Formatter: nested.Formatter{},
		})
		if e.cfg.Verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	})
	return e
}
