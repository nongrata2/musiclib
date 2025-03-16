FROM golang:1.23 AS build

COPY go.mod go.sum /src/
COPY api /src/api
COPY db /src/db

ENV CGO_ENABLED=0
RUN cd /src && go build -o /api api/main.go

FROM alpine:3.20

COPY --from=build /api /api

ENTRYPOINT [ "/api" ]