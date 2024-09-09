FROM golang:alpine as builder

ARG ARCH="arm64"

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
    
WORKDIR $GOPATH/src/svg-logos-downloader

COPY . .

RUN go mod download && go mod verify
RUN GOOS=linux GOARCH=${ARCH} go build -ldflags="-w -s" -o /go/bin/svg-logos-downloader cmd/api/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/svg-logos-downloader /go/bin/svg-logos-downloader

USER appuser:appuser

ENTRYPOINT ["/go/bin/svg-logos-downloader"]
