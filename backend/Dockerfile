FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o main .

FROM alpine:3.22

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/public/ /www/public/

# Config comes from environment variables at runtime (env_file in compose),
# so secrets never end up in the image.

EXPOSE 3000
ENTRYPOINT ["/www/main"]
