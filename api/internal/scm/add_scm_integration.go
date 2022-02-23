package scm

import (
	"errors"

	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type addSCMIntegrationRepository interface {
	Create(*SCMIntegration, *AccessToken) error
}

type AddSCMIntegrationDTO struct {
	Name string `json:"name"`
	Type SCMIntegrationType `json:"type"`
	Endpoint string `json:"endpoint"`
	AccessToken AddAccessTokenV1DTO `json:"access_token"`
}

type AddSCMIntegrationV1Response struct {
	Created *SCMIntegration
	errors []error
	status result.Status
}

func (r *AddSCMIntegrationV1Response) Data() interface {} {
	return *r.Created
}

func (r *AddSCMIntegrationV1Response) Meta() interface {} {
	return map[string]interface{}{}
}

func (r *AddSCMIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *AddSCMIntegrationV1Response) Status() result.Status {
	return r.status
}

func (r *AddSCMIntegrationV1Response) IsListData() bool {
	return false
}

func handleAddSCMIntegration(repo addSCMIntegrationRepository, dto AddSCMIntegrationDTO) result.Result {
	return &AddSCMIntegrationV1Response{
		Created: nil,
		errors: []error{
			errors.New("not implemented"),
		},
		status: result.InternalError,
	}
}

