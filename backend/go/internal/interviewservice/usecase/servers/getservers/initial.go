package getservers

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

func New(
	serverRepo repositorysql.ServerRepository,
	validator validatorx.Validator,
) *Handler {
	service := newService(serverRepo)
	handler := newHandler(service, validator)
	return handler
}
