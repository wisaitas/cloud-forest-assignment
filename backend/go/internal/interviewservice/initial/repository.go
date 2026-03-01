package initial

import (
	appRepositorySQL "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
)

type repository struct {
	userRepository        appRepositorySQL.UserRepository
	serverRepository      appRepositorySQL.ServerRepository
	activityLogRepository appRepositorySQL.ActivityLogRepository
}

func newRepository(sdk *sdk) *repository {
	userRepo, err := appRepositorySQL.NewUserRepository(sdk.bcrypt)
	if err != nil {
		panic(err)
	}
	return &repository{
		userRepository:        userRepo,
		serverRepository:      appRepositorySQL.NewServerRepository(),
		activityLogRepository: appRepositorySQL.NewActivityLogRepository(),
	}
}
