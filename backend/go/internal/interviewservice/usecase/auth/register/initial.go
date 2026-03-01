package register

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

func New(
	userRepository repositorysql.UserRepository,
	validator validatorx.Validator,
) *Handler {
	service := newService(userRepository)
	handler := newHandler(service, validator)

	return handler
}
