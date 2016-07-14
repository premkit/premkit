package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/premkit/premkit/log"
	"github.com/premkit/premkit/models"
)

// RegisterServiceParams contains parameters to the register service route.
// swagger:parameters registerService
type RegisterServiceParams struct {
	// Service registration parameters.
	// In: body
	Service         *models.Service `json:"service"`
	ReplaceExisting bool            `json:"replace_existing"`
}

// RegisterServiceResponse represents the response to a registerService call. This response
// includes a pointer to the registered service.
// swagger:response registerServiceResponse
type RegisterServiceResponse struct {
	// Service
	// In: body
	Body *models.Service `json:"service"`
}

// RegisterService is the handler called when a POST is made to register a new service.
func RegisterService(response http.ResponseWriter, request *http.Request) {
	// swagger:route POST /service services registerService
	//
	// Registers a new backend service with the router.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: https
	//
	//     Responses:
	//       201: registerServiceResponse
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	registerServiceParams := RegisterServiceParams{}
	if err := json.Unmarshal(body, &registerServiceParams); err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	service, err := registerService(&registerServiceParams)
	if err != nil {
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	registerServiceResponse := RegisterServiceResponse{
		Body: service,
	}
	b, err := json.Marshal(registerServiceResponse)
	if err != nil {
		log.Error(err)
		http.Error(response, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
	}

	response.WriteHeader(http.StatusCreated)
	response.Write(b)
}

func registerService(params *RegisterServiceParams) (*models.Service, error) {
	if params.ReplaceExisting {
		_, err := models.DeleteServiceByName([]byte(params.Service.Name))
		if err != nil {
			return nil, err
		}
	}

	service, err := models.CreateService(params.Service)
	if err != nil {
		return nil, err
	}

	return service, nil
}
