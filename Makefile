.PHONY: docs dapr minio
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

dapr:
	dapr run --app-id=video-store --dapr-http-port 3500 --dapr-grpc-port 50001 --components-path ./dapr/components

dapr-debug:
	dapr run --app-id=video-store --dapr-http-port 3500 --app-port 8081 --dapr-grpc-port 50001 --components-path ./dapr/components


# S3-like storage, used in dev
minio:
	docker run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"


# https://github.com/golang/mock
mockgen:
	mockgen -source ./internal/object-storage/object-storage.go -destination ./internal/mock/object-storage/object-storage.go
	mockgen -source ./internal/progress-broker/progress-broker.go -destination ./internal/mock/progress-broker/progress-broker.go
	mockgen -source ./internal/video-hosting/video-host.go -destination ./internal/mock/video-hosting/video-host.go
	mockgen github.com/dapr/go-sdk/client Client  > ./internal/mock/dapr/dapr-client.go