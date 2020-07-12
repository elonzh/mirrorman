package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/elonzh/mirrorman/pkg/server"
)

func (e *Executor) initServeCommand() {
	cmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			s := server.NewServer(e.cfg)
			s.Serve()
		},
	}
	cmd.Flags().StringP("addr", "a", ":9876", "proxy server address")
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		panic(err)
	}
	e.rootCmd.AddCommand(cmd)
}
