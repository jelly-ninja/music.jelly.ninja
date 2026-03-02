FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY ssh/go.mod ssh/go.sum ./
RUN go mod download

COPY ssh/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o strudel-wish

FROM alpine:3.19

RUN apk add --no-cache openssh-keygen

WORKDIR /app

COPY --from=builder /app/strudel-wish .
RUN mkdir -p .ssh && ssh-keygen -A && mv /etc/ssh/ssh_host_* .ssh/

EXPOSE 23234

CMD ["./strudel-wish", "--host", "0.0.0.0", "--port", "23234"]
