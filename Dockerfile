FROM golang:1.16 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY action action
RUN go build -o /out/action

FROM gcr.io/distroless/base@sha256:5e3fac1733c75e0e879a9770724e3960610a5cfbbfb5366559fbc334fe86c249

COPY --from=build /out/action /bin/action

ENTRYPOINT ["/bin/action"]
