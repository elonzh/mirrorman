package cmd

import (
	"github.com/spf13/cobra"

	"github.com/elonzh/mirrorman/pkg/server"
)

func NewServeCommand() *cobra.Command {
	var (
		verbose  bool
		httpAddr string
	)
	cmd := &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			s := server.NewServer()
			s.Addr = httpAddr
			s.Proxy.Verbose = verbose
			s.Init()
			s.Serve()
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", true, "should every proxy request be logged to stdout")
	cmd.Flags().StringVarP(&httpAddr, "addr", "a", ":9876", "proxy server address")
	return cmd
}
