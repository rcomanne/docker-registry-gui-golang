FROM golang:1.18-alpine AS builder

# Add git for go mod and ca-certificates for SSL certs
RUN apk update && \
    apk --no-cache --no-progress add ca-certificates git && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Go into the module's directory
WORKDIR $GOPATH/src/github.com/rcomanne/docker-registry-gui
# Copy all the local files
COPY . .
# Download the go modules
RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/docker-registry-gui && \
    chmod 0777 ./docker-registry-gui

# The actual image to run the application
FROM alpine:3

# Copy the certs
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
# Copy the binary
COPY --from=builder /go/src/github.com/rcomanne/docker-registry-gui/docker-registry-gui /docker-registry-gui

CMD ["/docker-registry-gui"]