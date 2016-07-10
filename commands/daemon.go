package commands

import (
	"github.com/premkit/premkit/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	port     int
	dataFile string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "A high performance reverse proxy",
	Long: `PremKit is built around this high performance reverse proxy which supports 
backend service automatic registration with namespacing.  

'premkit daemon' will start this reverse proxy and listen for requests for registrations 
and for requests to forward to backend services.`,
}

func init() {
	daemonCmd.Flags().IntVarP(&port, "bind", "b", 80, "port on which the reverse proxy will bind and listen")
	daemonCmd.Flags().StringVarP(&dataFile, "datafile", "d", "/data/premkit.db", "location of the database file")

	daemonCmd.RunE = daemon
}

func daemon(cmd *cobra.Command, args []string) error {
	viper.Set("data_file", dataFile)
	server.Run(80)

	return nil
}
