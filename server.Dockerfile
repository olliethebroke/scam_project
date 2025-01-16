FROM golang:1.23.4-alpine as build-stage
COPY .  /server/
WORKDIR /server/
RUN go mod download
RUN go build -o ./bin/server_out ./cmd/app/server.go

FROM alpine:latest
ARG ENV_FILE
COPY --from=build-stage /server/bin/server_out .
ADD ${ENV_FILE} .
CMD ["./server_out"]