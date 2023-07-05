## Build
FROM golang:latest AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.11

COPY . .

ENV GOOS linux
ENV CGO_ENABLED 0

RUN make swagger
RUN make build

## Deploy
FROM alpine:3.18.2

WORKDIR /

COPY --from=build /app/build/server /server

ENV GIN_MODE=release

ENTRYPOINT ["/server"]

