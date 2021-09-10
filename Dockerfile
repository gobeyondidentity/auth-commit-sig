FROM golang:1.16 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY action action

ARG VERSION=latest
RUN go build \
  -ldflags "-X byndid/auth-commit-sig/action.version=${VERSION}" \
  -o /out/action

FROM debian:11.0-slim

COPY --from=build /out/action /bin/action

ENTRYPOINT ["/bin/action"]
