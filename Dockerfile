# FROM golang:latest as builder
FROM golang:1.24.4-alpine3.21
COPY . /auth
WORKDIR /auth
RUN go mod download && go build -o server .
EXPOSE 3000
CMD ["./server"]

# FROM scratch
# WORKDIR /app
# COPY --from=builder /app/publish .
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert
# ENV GIN_MODE=release \
#     PORT=80
# EXPOSE 80
# ENTRYPOINT ["./toc-generator"]