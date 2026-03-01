.PHONY: run-infra-service run schema

run-infra-service:
	cd backend/go/cmd/infraservice && go run main.go

run:
	cd backend/go && go run ./cmd/interviewservice

gen:
	cd backend/go && go generate ./...