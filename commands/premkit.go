package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PremkitCmd is the main (root) command for the CLI.
var PremkitCmd = &cobra.Command{
	Use:   "premkit",
	Short: "premkit enables installable software",
	Long:  "premkit is the reverse proxy, used to ship installable software",

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InitializeConfig(); err != nil {
			return err
		}

		cmd.Usage()
		os.Exit(-1)
		return nil
	},
}

// Execute adds all child commands to the room command PremkitCmd and sets flags
func Execute() {
	AddCommands()

	if _, err := PremkitCmd.ExecuteC(); err != nil {
		os.Exit(-1)
	}
}

// AddCommands will add all child commands to the PremkitCmd
func AddCommands() {
	PremkitCmd.AddCommand(daemonCmd)
}

// InitializeConfig initializes the config environment with defaults.
func InitializeConfig(subCmdVs ...*cobra.Command) error {
	viper.SetEnvPrefix("premkit")
	return nil
}
