package commands

import (
	"fmt"

	"github.com/premkit/premkit/certs"
	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = 80
	defaultHTTPSPort              = 443
	defaultTLSKeyFile             = ""
	defaultTLSCertFile            = ""
	defaultGenerateSelfSignedCert = true

	defaultDataFile = "/data/premkit.db"
	defaultTLSStore = "/data/premkit.db"
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
	viper.SetEnvPrefix("premkit")
	viper.AutomaticEnv()

	daemonCmd.Flags().Int("bind-http", defaultHTTPPort, "port on which the reverse proxy will bind and listen for http connections")
	daemonCmd.Flags().Int("bind-https", defaultHTTPSPort, "port on which the reverse proxy will bind and listen for https (tls) connections")
	daemonCmd.Flags().String("key-file", defaultTLSKeyFile, "path to private key to use when serving tls connections")
	daemonCmd.Flags().String("cert-file", defaultTLSCertFile, "path to cert to use when serving tls connections")
	daemonCmd.Flags().Bool("self-signed", defaultGenerateSelfSignedCert, "true to have premkit generate a new self-signed keypair to use for tls connections")
	daemonCmd.Flags().String("data-file", defaultDataFile, "location of the database file")
	daemonCmd.Flags().String("tls-store", defaultTLSStore, "location to store generated tls certs and keys in")

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

	showAppliedSettings()

	server.Run(config)
	return nil
}

func buildConfig() (*server.Config, error) {
	keyFile := ""
	certFile := ""

	if viper.GetInt("bind_https") != 0 {
		if viper.GetBool("self_signed") {
			k, c, err := certs.GenerateSelfSigned(viper.GetString("tls_store"))
			if err != nil {
				log.Error(err)
				return nil, err
			}
			keyFile = k
			certFile = c
		} else if viper.GetString("key_file") != "" && viper.GetString("cert_file") != "" {
			if err := certs.ParseKeyPair(viper.GetString("key_file"), viper.GetString("cert_file")); err != nil {
				log.Error(err)
				return nil, err
			}

			keyFile = viper.GetString("key_file")
			certFile = viper.GetString("cert_file")
		}
	}

	config := server.Config{
		HTTPPort:  viper.GetInt("bind_http"),
		HTTPSPort: viper.GetInt("bind_https"),

		TLSKeyFile:  keyFile,
		TLSCertFile: certFile,
	}

	return &config, nil
}

func showAppliedSettings() {
	var nonDefault []string

	if viper.GetInt("bind_http") != defaultHTTPPort {
		nonDefault = append(nonDefault, fmt.Sprintf("HTTP Bind Port set to %d", viper.GetInt("bind_http")))
	}
	if viper.GetInt("bind_https") != defaultHTTPSPort {
		nonDefault = append(nonDefault, fmt.Sprintf("HTTPS Bind Port set to %d", viper.GetInt("bind_https")))
	}

	if viper.GetString("key_file") != defaultTLSKeyFile {
		nonDefault = append(nonDefault, fmt.Sprintf("TLS Key File set to %s", viper.GetString("key_file")))
	}
	if viper.GetString("cert_file") != defaultTLSCertFile {
		nonDefault = append(nonDefault, fmt.Sprintf("TLS Cert File set to %s", viper.GetString("cert_file")))
	}
	if viper.GetBool("self_signed") != defaultGenerateSelfSignedCert {
		nonDefault = append(nonDefault, fmt.Sprintf("Generate Self Signed Cert set to %v", viper.GetBool("self_signed")))
	}

	if viper.GetString("data_file") != defaultDataFile {
		nonDefault = append(nonDefault, fmt.Sprintf("DataFile set to %s", viper.GetString("data_file")))
	}
	if viper.GetString("tls_store") != defaultTLSStore {
		nonDefault = append(nonDefault, fmt.Sprintf("TLS Store set to %s", viper.GetString("tls_store")))
	}

	if len(nonDefault) == 0 {
		log.Infof("Using default settings")
		return
	}

	log.Infof("Overridden settings: ")
	for _, n := range nonDefault {
		log.Infof("\t%s", n)
	}
}
