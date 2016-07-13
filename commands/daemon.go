package commands

import (
	"github.com/premkit/premkit/certs"
	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	httpPort               int
	httpsPort              int
	tlsKeyFile             string
	tlsCertFile            string
	generateSelfSignedCert bool

	dataFile string
	tlsStore string
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
	daemonCmd.Flags().IntVarP(&httpPort, "bind-http", "", 80, "port on which the reverse proxy will bind and listen for http connections")
	daemonCmd.Flags().IntVarP(&httpsPort, "bind-https", "", 443, "port on which the reverse proxy will bind and listen for https (tls) connections")
	daemonCmd.Flags().StringVarP(&tlsKeyFile, "key-file", "", "", "path to private key to use when serving tls connections")
	daemonCmd.Flags().StringVarP(&tlsCertFile, "cert-file", "", "", "path to cert to use when serving tls connections")
	daemonCmd.Flags().BoolVarP(&generateSelfSignedCert, "self-signed", "", true, "true to have premkit generate a new self-signed keypair to use for tls connections")
	daemonCmd.Flags().StringVarP(&dataFile, "data-file", "", "/data/premkit.db", "location of the database file")
	daemonCmd.Flags().StringVarP(&tlsStore, "tls-store", "", "/data/tls", "location to store generated tls certs and keys in")

	viper.BindPFlag("bind_http", daemonCmd.Flags().Lookup("bind-http"))
	viper.BindPFlag("bind_https", daemonCmd.Flags().Lookup("bind-https"))
	viper.BindPFlag("key_file", daemonCmd.Flags().Lookup("key-file"))
	viper.BindPFlag("cert_file", daemonCmd.Flags().Lookup("cert-file"))
	viper.BindPFlag("self_signed", daemonCmd.Flags().Lookup("self-signed"))
	viper.BindPFlag("data_file", daemonCmd.Flags().Lookup("data-file"))
	viper.BindPFlag("tls_store", daemonCmd.Flags().Lookup("tls-store"))

	daemonCmd.RunE = daemon
}

func daemon(cmd *cobra.Command, args []string) error {
	if err := InitializeConfig(daemonCmd); err != nil {
		return err
	}

	config, err := buildConfig()
	if err != nil {
		return err
	}

	server.Run(config)
	return nil
}

func buildConfig() (*server.Config, error) {
	keyFile := ""
	certFile := ""

	if httpsPort != 0 {
		if generateSelfSignedCert {
			k, c, err := certs.GenerateSelfSigned(tlsStore)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			keyFile = k
			certFile = c
		} else if tlsKeyFile != "" && tlsCertFile != "" {
			if err := certs.ParseKeyPair(tlsKeyFile, tlsCertFile); err != nil {
				log.Error(err)
				return nil, err
			}

			keyFile = tlsKeyFile
			certFile = tlsCertFile
		}
	}

	config := server.Config{
		HTTPPort:  httpPort,
		HTTPSPort: httpsPort,

		TLSKeyFile:  keyFile,
		TLSCertFile: certFile,
	}

	return &config, nil
}
