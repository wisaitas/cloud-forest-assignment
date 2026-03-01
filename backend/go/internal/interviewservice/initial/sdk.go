package initial

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"
	"github.com/wisaitas/cloud-forest-assignment/pkg/bcryptx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/jwtx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

type sdk struct {
	validator validatorx.Validator
	bcrypt    bcryptx.Bcrypt
	jwt       jwtx.Jwt
}

func newSDK() *sdk {
	return &sdk{
		validator: validatorx.NewValidator(),
		bcrypt:    bcryptx.NewBcrypt(),
		jwt:       jwtx.NewJwt(interviewservice.Config.Jwt.Secret),
	}
}
