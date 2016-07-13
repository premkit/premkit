package server

// Config represents the config to use to start the web server.
type Config struct {
	HTTPPort  int
	HTTPSPort int

	TLSKeyFile  string
	TLSCertFile string
}
