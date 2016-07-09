package commands

import (
	"os"

	"github.com/spf13/cobra"
)

// PremKitCmd is the main (root) command for the CLI.
var PremKitCmd = &cobra.Command{
	Use:   "premkit",
	Short: "premkit enables installable software",
	Long:  "premkit is the reverse proxy, used to ship installable software",

	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Usage()
		os.Exit(-1)
		return nil
	},
}

// Execute adds all child commands to the room command PremKitCmd and sets flags
func Execute() {
	AddCommands()

	if _, err := PremKitCmd.ExecuteC(); err != nil {
		os.Exit(-1)
	}
}

// AddCommands will add all child commands to the PremKitCmd
func AddCommands() {
	PremKitCmd.AddCommand(daemonCmd)
}
