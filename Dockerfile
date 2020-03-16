# Builder image
FROM golang:1.13-alpine as builder

# Copy source
ADD ./cmd/kubenab /root/app

# Install dependencies
RUN apk add --no-cache git

# Download modules
RUN ls -alh /root/app && cd /root/app && \
    GO111MODULE=on GOPROXY=https://gocenter.io go mod download

# Build microservices
RUN cd /root/app && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM gcr.io/distroless/static

# Copy microservice executable from builder image
COPY --from=builder /root/app/kubenab /bin/kubenab

# Set runtime user to non-root
USER 1000

# Set Entrypoint
CMD ["kubenab"]
