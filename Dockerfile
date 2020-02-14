# Builder image
FROM golang:1.12-alpine as builder

# Install dependencies
RUN apk add --no-cache git

# Set workspace
WORKDIR /src/kubenab/kubenab/

# Copy only go.sum and go.mod to allow for docker cache of modules
COPY cmd/kubenab/go.mod cmd/kubenab/go.sum /src/kubenab/kubenab/cmd/kubenab/

# Download modules
RUN cd cmd/kubenab && \
    GO111MODULE=on GOPROXY=https://gocenter.io go mod download

# Copy source
COPY . /src/kubenab/kubenab/

# Build microservices
RUN apk add make git && \
    OUT_DIR=/ make build

FROM gcr.io/distroless/static

# Copy microservice executable from builder image
COPY --from=builder /kubenab /bin/kubenab

# Set runtime user to non-root
USER 1000

# Set Entrypoint
CMD ["kubenab"]
