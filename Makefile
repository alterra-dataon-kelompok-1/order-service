.PHONY: run
run: ## run the API server in development environment
	go run main.go

.PHONY: dev
dev: ## run the API server in development environment
	go run main.go -env=dev

.PHONY: test
test: ## test with coverage status
	go test ./... -cover

.PHONY: test-cover-report
test-cover-report: ## test, generate coverage report, and show it
	go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

.PHONY: mock-repository
mock-repository: ## generate mock repository
	mockgen -source=./internal/repository/repository.go \
	-destination=./internal/app/order/mock/repository_mock.go

.PHONY: mock-service
mock-service: ## test, generate coverage report, and show it
	mockgen -source=./internal/app/order/service.go \
	-destination=./internal/app/order/mock/service_mock.go