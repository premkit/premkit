package networking

import (
	"crypto/tls"
	"os"
)

func GetTLSConfig(certs []tls.Certificate) (tlsConfig *tls.Config) {
	if os.Getenv("USE_STRICT_TLS") == "false" {
		tlsConfig = &tls.Config{
			Certificates: certs,
		}
	} else {
		tlsConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: certs,
		}
	}

	return
}
