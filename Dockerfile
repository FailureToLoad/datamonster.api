FROM golang:1.21.11-alpine3.20 AS build-stage

WORKDIR /app

COPY ./ ./
RUN go mod download
RUN go build -o /out/apiserver ./cmd/apiserver/main.go
RUN ls /out
FROM gcr.io/distroless/base-debian12:latest AS run-stage

WORKDIR /


COPY --from=build-stage /out/apiserver /apiserver

USER nonroot:nonroot
ENTRYPOINT ["./apiserver"]