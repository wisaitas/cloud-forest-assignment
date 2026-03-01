package initial

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/caller/infraservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/auth"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers"
)

var MockInMemmory = map[string]any{}

type useCase struct {
	authUseCase    *auth.UseCase
	serversUseCase *servers.UseCase
}

func newUseCase(
	config *config,
	repository *repository,
	sdk *sdk,
) *useCase {
	infraClient := infraservice.NewClient(interviewservice.Config.InfraService.URL)
	return &useCase{
		authUseCase:    auth.New(repository.userRepository, sdk.bcrypt, sdk.jwt, sdk.validator),
		serversUseCase: servers.New(repository.serverRepository, repository.activityLogRepository, infraClient, sdk.validator),
	}
}
