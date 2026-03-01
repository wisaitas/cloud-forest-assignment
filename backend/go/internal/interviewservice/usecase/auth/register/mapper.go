package register

import "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"

func (s *service) mapRequestToEntity(request *Request) *entity.User {
	return &entity.User{
		Email:    request.Email,
		Password: request.Password,
	}
}
