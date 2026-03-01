package provision

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/caller/infraservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

func New(
	serverRepo repositorysql.ServerRepository,
	activityLogRepo repositorysql.ActivityLogRepository,
	infraClient infraservice.InfraServiceCaller,
	validator validatorx.Validator,
) *Handler {
	service := newService(serverRepo, activityLogRepo, infraClient)
	handler := newHandler(service, validator)
	return handler
}
