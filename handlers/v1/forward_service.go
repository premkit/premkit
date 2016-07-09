package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/models"

	"github.com/vulcand/oxy/forward"
)

var (
	fwd *forward.Forwarder
)

func init() {
	f, err := forward.New()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	fwd = f
}

// TODO Swagger this handler, but it's a special handler and therefore a little trickier to swagger-ify

// ForwardService is the handler for anything that should be possibly fowarded to an upstream.
func ForwardService(response http.ResponseWriter, request *http.Request) {
	// TODO keep these cached because in any reasonable load this will be painful
	services, err := models.ListServices()
	if err != nil {
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	var url *url.URL

	for _, service := range services {
		path := strings.TrimPrefix(request.URL.Path, "/")

		pathParts := strings.Split(path, "/")
		if len(pathParts) == 0 {
			continue
		}

		if pathParts[0] != service.Path {
			continue
		}

		log.Debugf("path %q matched service %q (service path %q)", request.URL.Path, service.Name, service.Path)

		if len(service.Upstreams) == 0 {
			err := errors.New("No upstreams are available")
			log.Error(err)
			response.WriteHeader(http.StatusBadGateway)
			response.Write([]byte(""))
			return
		}

		childPath := strings.Join(pathParts[1:], "/")
		request.RequestURI = "/" + childPath

		log.Debugf("service.Upstreams = %s/%s", service.Upstreams[0], childPath)
		u, err := url.Parse(fmt.Sprintf("%s/%s", service.Upstreams[0], childPath))
		if err != nil {
			log.Error(err)
			http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
			return
		}

		url = u
	}

	if url == nil {
		http.Error(response, "Route not found", http.StatusNotFound)
		return
	}

	request.URL = url
	fwd.ServeHTTP(response, request)
}
