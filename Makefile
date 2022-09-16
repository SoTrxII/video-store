.PHONY: docs
docs:
# To install https://github.com/swaggo/gin-swagger
	swag init

createBinDir:
# "true" command is no universal so we'll use echo instead
	mkdir bin || echo "true"

build: createBinDir
	go build -o bin/server

run : build
	./bin/server

release: createBinDir
	go build -ldflags "-s -w" -o bin/server

test:
	go test $(shell go list ./... | grep -v /mock/ | grep -v /docs | grep -v /logger | grep -v /youtube-host | grep -v /video-host)  -covermode=atomic -coverprofile=coverage.out


# https://github.com/golang/mock
mockgen:
	mockgen -source ./internal/object-storage/object-storage.go -destination ./internal/mock/object-storage/object-storage.go
	mockgen -source ./internal/progress-broker/progress-broker.go -destination ./internal/mock/progress-broker/progress-broker.go
	mockgen -source ./internal/video-hosting/video-host.go -destination ./internal/mock/video-hosting/video-host.go
	mockgen github.com/dapr/go-sdk/client Client  > ./internal/mock/dapr/dapr-client.go