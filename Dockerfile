FROM golang:1.18-alpine AS builder

RUN apk --no-cache --no-progress add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Go into the module's directory
WORKDIR /go/src/github.com/rcomanne/docker-registry-gui
# Download the go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the app
COPY cmd/ /go/src/github.com/rcomanne/docker-registry-gui/cmd/
COPY pkg/ /go/src/github.com/rcomanne/docker-registry-gui/pkg/

RUN go build ./cmd/docker-registry-gui && \
    chmod a+x docker-registry-gui

FROM alpine:3 AS runner

RUN id && \
    cat /etc/passwd

WORKDIR /
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/rcomanne/docker-registry-gui/docker-registry-gui /docker-registry-gui

ENTRYPOINT ["/docker-registry-gui"]