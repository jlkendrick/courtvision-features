FROM golang:1.22.4-alpine

WORKDIR /app

COPY ./ .

EXPOSE 8080

CMD ["go", "run", "v2.go"]

# Build command: docker build -t stopz-server .
# Run command: docker run -p 8000:8000 stopz-server