# Builder image
FROM golang:1.12-alpine as builder

# Install dependencies
RUN apk add --no-cache git

# Set workspace
WORKDIR /src/kubenab/kubenab/

# Copy source
COPY ./ /src/kubenab/kubenab/

# Download modules
RUN cd cmd/kubenab && \
    GO111MODULE=on GOPROXY=https://gocenter.io go mod download

# Build microservices
RUN apk add make git && \
    OUT_DIR=/ make build

FROM gcr.io/distroless/static

# Copy microservice executable from builder image
COPY --from=builder /kubenab /bin/kubenab

# Set Entrypoint
CMD ["kubenab"]
