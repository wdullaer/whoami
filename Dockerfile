FROM golang:1.12-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0
RUN apk --update add git

# Copy dependencies first to improve caching
COPY go.mod .
COPY go.sum .
COPY vendor/ vendor/

COPY . .

# Build the service and inject compile time variables.
RUN \
go build -mod=vendor

FROM alpine:latest

CMD ["/whoami"]
EXPOSE 50501
RUN apk --update add ca-certificates
COPY --from=builder /app/whoami /whoami
