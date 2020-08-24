package commands

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

type jsonVersion struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func (e *Executor) initVersionCommand() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ver := jsonVersion{
				Version: e.version,
				Commit:  e.commit,
				Date:    e.date,
			}
			data, err := json.Marshal(&ver)
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		},
	}
	e.rootCmd.AddCommand(versionCmd)
}
