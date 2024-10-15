FROM golang:1.22.4 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o exec

FROM alpine:latest AS certs

RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /app

COPY --from=builder /app/exec /app/exec

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./static/schedule24-25.json /app/static/schedule24-25.json

CMD ["./exec"]

# Build command: docker build -t stopz-server .
# Run command: docker run -p 8000:8000 stopz-server