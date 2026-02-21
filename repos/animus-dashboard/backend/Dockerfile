FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download || true

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o animus-api ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates ansible openssh-client

WORKDIR /app

COPY --from=builder /app/animus-api .
COPY ansible/ ./ansible/

RUN adduser -D -u 1000 animus
RUN chown -R animus:animus /app
USER animus

EXPOSE 8080

CMD ["./animus-api"]
