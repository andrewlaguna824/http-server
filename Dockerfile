# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build_stage
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build cmd/main.go

FROM alpine
COPY --from=build_stage /app/main .
# COPY --from=build_stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
CMD ./main
