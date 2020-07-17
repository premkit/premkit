package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	v1 "github.com/premkit/premkit/handlers/v1"
	"github.com/premkit/premkit/log"

	"github.com/replicatedcom/replicated/pkg/networking"

	"github.com/gorilla/mux"
)

// Run is the main entrypoint of this daemon.
func Run(config *Config) error {
	router := mux.NewRouter()

	internal := router.PathPrefix("/premkit").Subrouter()
	internalV1 := internal.PathPrefix("/v1").Subrouter()
	internalV1.HandleFunc("/service", v1.RegisterService).Methods("POST")

	// TODO serve the swagger.json using a gorilla static handlers

	forward := router.PathPrefix("/").Subrouter()
	forward.HandleFunc("/{path:.*}", v1.ForwardService)

	if config.HTTPPort != 0 {
		go func() {
			log.Infof("Listening on port %d for http connections", config.HTTPPort)
			log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPPort), router))
		}()
	}

	pair, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
	if err != nil {
		log.Errorf("Failed to load x509 key pair: %v", err)
		return err
	}

	if config.HTTPSPort != 0 {
		go func() {
			log.Infof("Listening on port %d for https connections", config.HTTPSPort)
			srv := &http.Server{
				Addr:      fmt.Sprintf(":%d", config.HTTPSPort),
				Handler:   router,
				TLSConfig: networking.GetTLSConfig([]tls.Certificate{pair}),
			}
			log.Error(srv.ListenAndServeTLS("", ""))
		}()
	}

	<-make(chan struct{})
	return nil
}
