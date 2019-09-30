# First stage: build the executable
FROM golang:1 AS builder

# Enable Go modules and vendor
ENV GO111MODULE=on GOFLAGS=-mod=vendor

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd cmd
COPY go.mod go.sum main.go ./
COPY vendor vendor
COPY apply apply
COPY config config
COPY docs docs
COPY errs errs
COPY exp exp
COPY init init
COPY migrations migrations
COPY plan plan
COPY plugins plugins
COPY setup setup
COPY templates templates
COPY testdata testdata
COPY util util

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o fogg .

# Install chamber
ENV CHAMBER_VERSION=v2.1.0
RUN wget -q https://github.com//segmentio/chamber/releases/download/${CHAMBER_VERSION}/chamber-${CHAMBER_VERSION}-linux-amd64 -O /bin/chamber
RUN chmod +x /bin/chamber

# Final stage: the running container
FROM alpine:latest AS final

# Install SSL root certificates
RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /app/fogg /bin/fogg
COPY --from=builder /bin/chamber /bin/chamber

CMD ["fogg"]