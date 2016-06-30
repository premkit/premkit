package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/models"
	"github.com/premkit/premkit/schema"

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

// RegisterService is the handler called when a POST is made to register a new service.
func RegisterService(response http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	registerServiceRequest := schema.RegisterServiceRequest{}
	if err := json.Unmarshal(body, &registerServiceRequest); err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	service, err := models.CreateService(registerServiceRequest.Service)
	if err != nil {
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	registerServiceResponse := schema.RegisterServiceResponse{
		Service: service,
	}
	b, err := json.Marshal(registerServiceResponse)
	if err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
	}

	response.WriteHeader(http.StatusCreated)
	response.Write(b)
}

// MaybeForwardService is the handler for anything that should be possibly fowarded to an upstream.
func MaybeForwardService(response http.ResponseWriter, request *http.Request) {
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
