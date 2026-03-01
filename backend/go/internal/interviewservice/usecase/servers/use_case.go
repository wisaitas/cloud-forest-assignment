package servers

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/caller/infraservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/getservers"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/power"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/provision"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

type UseCase struct {
	List      *getservers.Handler
	Provision *provision.Handler
	Power     *power.Handler
}

func New(
	serverRepo repositorysql.ServerRepository,
	activityLogRepo repositorysql.ActivityLogRepository,
	infraClient infraservice.InfraServiceCaller,
	validator validatorx.Validator,
) *UseCase {
	return &UseCase{
		List:      getservers.New(serverRepo, validator),
		Provision: provision.New(serverRepo, activityLogRepo, infraClient, validator),
		Power:     power.New(serverRepo, activityLogRepo, infraClient, validator),
	}
}
