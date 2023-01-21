FROM golang:1.17.2-alpine AS builder
RUN apk update && apk add --no-cache git make openssh-client curl ca-certificates
RUN update-ca-certificates
RUN curl https://repos.balad.ir/artifactory/github/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.2/grpc_health_probe-linux-amd64 -o /bin/grpc_health_probe && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app
COPY Makefile go.mod go.sum ./
RUN make download-dependencies

COPY . .
RUN make build
RUN chmod +x dvs


FROM scratch

COPY --from=builder /app/dvs /app/dvs
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENV USER root
ENV HOME /app
CMD ["./dvs", "serve"]
