FROM golang:1.16 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY action action
RUN go build -o /out/action

FROM gcr.io/distroless/base

COPY --from=build /out/action /bin/action

ENTRYPOINT ["/bin/action"]
