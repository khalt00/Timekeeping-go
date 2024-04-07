FROM golang:1.22-alpine3.19 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

FROM alpine:3.19
WORKDIR /
COPY --from=builder /app/server /app/server
COPY --from=builder /app/config.yaml /config.yaml
COPY --from=builder /app/abi /abi
COPY --from=builder /app/db/migrations /db/migrations

EXPOSE 8080
CMD ["/app/server"]
