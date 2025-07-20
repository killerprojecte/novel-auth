FROM golang:1.24.5-alpine3.22 AS builder
WORKDIR /api
COPY . .
RUN go build -o api-server

FROM alpine:3.22.1
COPY --from=builder /api/api-server /api/api-server
WORKDIR /api
EXPOSE 3000
CMD ["./api-server"]