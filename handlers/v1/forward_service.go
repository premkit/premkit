package v1

import (
	"crypto/tls"
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
	fwdSecure   *forward.Forwarder
	fwdInsecure *forward.Forwarder
)

func init() {
	var insecureTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	insecureRoundTripper := forward.RoundTripper(insecureTransport)
	f, err := forward.New(insecureRoundTripper)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	fwdInsecure = f

	var secureTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	secureRoundTripper := forward.RoundTripper(secureTransport)
	f, err = forward.New(secureRoundTripper)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	fwdSecure = f
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

	var service *models.Service

	log.Debugf("Looking for a known route with prefix %q", request.URL.Path)
	for _, s := range services {
		if !isPathPrefix(s.Path, request.URL.Path) {
			log.Debugf("Service with path %q did not match", s.Path)
			continue
		}

		log.Debugf("path %q matched service %q (service path %q)", request.URL.Path, s.Name, s.Path)
		service = s

		break
	}

	if service == nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(""))
		return
	}

	if len(service.Upstreams) == 0 {
		err := errors.New("No upstreams are available")
		log.Error(err)
		response.WriteHeader(http.StatusBadGateway)
		response.Write([]byte(""))
		return
	}

	// TODO pick an upstream with some intelligence

	// The upstream we will forward to
	upstream := service.Upstreams[0]

	url, err := getForwardURLForServiceRequest(upstream, service, request.URL.Path)
	if err != nil {
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	if url == nil {
		http.Error(response, "Route not found", http.StatusNotFound)
		return
	}

	request.URL = url
	request.RequestURI = url.Path

	if upstream.InsecureSkipVerify {
		fwdInsecure.ServeHTTP(response, request)
	}

	fwdSecure.ServeHTTP(response, request)
}

func stripLeadingSlashIfPresent(path string) string {
	return strings.TrimPrefix(path, "/")
}

func isPathPrefix(servicePath, requestPath string) bool {
	servicePath = stripLeadingSlashIfPresent(servicePath)
	requestPath = stripLeadingSlashIfPresent(requestPath)

	return strings.HasPrefix(requestPath, servicePath)
}

func createForwardPath(servicePath, requestPath string) string {
	servicePath = stripLeadingSlashIfPresent(servicePath)
	requestPath = stripLeadingSlashIfPresent(requestPath)

	// Remove the servicePath from the requestPath
	return strings.TrimPrefix(requestPath, servicePath)
}

// TODO refactor this out into a new module
func getForwardURLForServiceRequest(upstream *models.Upstream, service *models.Service, path string) (*url.URL, error) {
	childPath := createForwardPath(service.Path, path)
	//request.RequestURI = childPath

	// The built url we will forward to
	upstreamURL := ""
	if upstream.IncludeServicePath {
		upstreamURL = fmt.Sprintf("%s/%s%s", upstream.URL, service.Path, childPath)
	} else {
		upstreamURL = fmt.Sprintf("%s%s", upstream.URL, childPath)
	}

	upstreamURL = strings.TrimPrefix(upstreamURL, "/")

	url, err := url.Parse(upstreamURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return url, nil
}
