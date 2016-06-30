package schema

import (
	"github.com/premkit/premkit/models"
)

type RegisterServiceRequest struct {
	Service *models.Service `json:"service"`
}

type RegisterServiceResponse struct {
	Service *models.Service `json:"service"`
}
