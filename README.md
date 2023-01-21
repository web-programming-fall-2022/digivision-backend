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
