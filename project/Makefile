AUTH_BINARY=authApp



up:
	@echo starting Docker images...
	docker-compose up
	@echo docker images started



up_build: build_auth
	@echo Stopping docker images (if running...)
	docker-compose down

	@echo Building (when required) and starting docker images...
	docker-compose up  --build -d
	@echo Docker images built and started

down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

build_auth:
	@echo Building auth binary
	chdir ..\auth && ser GOOS=linux&& set GOARCH=amd64&& set
	CGO_ENABLED=0 && go build -o ${AUTH_BINARY} ./cmd/api
	@echo Done!