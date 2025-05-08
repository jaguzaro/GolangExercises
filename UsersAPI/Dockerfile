FROM golang:1.24-alpine AS builder
WORKDIR /build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app .

FROM alpine:3.17 AS final

COPY --from=builder /app /bin/app
EXPOSE 8080
CMD ["bin/app"]
