# Digivision Backend

A service to search and retrieve products from Digikala.

## Development

To download dependencies:
```shell script
make download-dependencies
```
To build:
```shell script
make build
```
To run unit tests:
```shell script
make test
```
To run locally after build:
```shell script
./dvs serve --dev
```

## Building protobufs

```shell script
protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=pkg/api/v1 --go-grpc_out=pkg/api/v1 --grpc-gateway_out=logtostderr=true:pkg/api/v1 --swagger_out=logtostderr=true:api/swagger/v1 search.proto
protoc --proto_path=api/proto/img2vec --go_out=internal/api/img2vec --go-grpc_out=internal/api/img2vec img2vec.proto
protoc --proto_path=api/proto/od --go_out=internal/api/od --go-grpc_out=internal/api/od object_detector.proto
```

## Build docker image
To build:
```shell script
make docker-build VERSION=0.0.1
```
To push to registry:
```shell script
make docker-push VERSION=0.0.1
```

## Deployment
## Local (docker compose)
if you want to deploy on development environment with docker compose:
```shell script
cp .env.sample .env
make deploy-dev
```
