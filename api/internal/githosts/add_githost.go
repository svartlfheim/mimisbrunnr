package githosts

import (
	"errors"

	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type AddGitHostDTO struct {
	Name *string
	Type *GitHostType
	Endpoint *string
	Credentials *AddCredentialsV1DTO
}

type AddGitHostV1Response struct {
	Created *GitHost
	errors []error
	status result.Status
}

func (r *AddGitHostV1Response) Data() interface {} {
	return *r.Created
}

func (r *AddGitHostV1Response) Meta() interface {} {
	return map[string]interface{}{}
}

func (r *AddGitHostV1Response) Errors() []error {
	return r.errors
}

func (r *AddGitHostV1Response) Status() result.Status {
	return r.status
}

func (r *AddGitHostV1Response) IsListData() bool {
	return false
}

func (m *Manager) Add(dto AddGitHostDTO) (result.Result) {
	return &AddGitHostV1Response{
		Created: nil,
		errors: []error{
			errors.New("not implemented"),
		},
		status: result.InternalError,
	}
}