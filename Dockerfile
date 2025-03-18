FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/musiclib ./cmd/musiclib
COPY internal ./internal
COPY migrations ./migrations

RUN CGO_ENABLED=0 go build -o /musiclib ./cmd/musiclib/main.go

FROM alpine:3.20

COPY --from=build /musiclib /musiclib

ENTRYPOINT ["/musiclib"]