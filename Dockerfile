# Fetch
FROM golang:latest AS fetch-stage
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

# Generate
FROM ghcr.io/a-h/templ:latest AS generate-stage
COPY --chown=65532:65532 . /app
WORKDIR /app
RUN ["templ", "generate"]

# Build
FROM golang:latest AS build-stage
COPY --from=generate-stage /app /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -ldflags="-w -s" -o /app/app cmd/api/main.go
RUN mkdir -p /uploads && chown 65532:65532 /uploads && chmod 755 /uploads

# Deploy
FROM gcr.io/distroless/base-debian12 AS deploy-stage
WORKDIR /
COPY --from=build-stage /app/app /app
COPY --from=build-stage --chown=65532:65532 /uploads /uploads
EXPOSE 8083
USER nonroot:nonroot
ENTRYPOINT ["/app"]
