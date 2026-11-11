FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /src/edge ./cmd/edge
RUN go build -o /src/migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /src/edge ./edge
COPY --from=builder /src/migrations ./migrations
COPY --from=builder /src/migrate ./migrate

ENTRYPOINT ["./edge"]