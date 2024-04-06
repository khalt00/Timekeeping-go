FROM golang:1.21.6-alpine3.19 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:3.19
WORKDIR /
COPY --from=builder /app/server /app/server
# COPY --from=builder /app/config.docker.yaml config.yaml

EXPOSE 8080
CMD ["/app/server"]
