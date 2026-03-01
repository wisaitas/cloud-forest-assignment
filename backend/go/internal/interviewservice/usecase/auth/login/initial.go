package login

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/bcryptx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/jwtx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

func New(
	userRepository repositorysql.UserRepository,
	bcrypt bcryptx.Bcrypt,
	jwt jwtx.Jwt,
	validator validatorx.Validator,
) *Handler {
	service := newService(userRepository, bcrypt, jwt)
	handler := newHandler(service, validator)

	return handler
}
