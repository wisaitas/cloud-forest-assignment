package auth

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/auth/login"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/auth/register"
	"github.com/wisaitas/cloud-forest-assignment/pkg/bcryptx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/jwtx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

type UseCase struct {
	Register *register.Handler
	Login    *login.Handler
}

func New(
	userRepository repositorysql.UserRepository,
	bcrypt bcryptx.Bcrypt,
	jwt jwtx.Jwt,
	validator validatorx.Validator,
) *UseCase {
	return &UseCase{
		Register: register.New(userRepository, validator),
		Login:    login.New(userRepository, bcrypt, jwt, validator),
	}
}
