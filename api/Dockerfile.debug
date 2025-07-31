FROM golang:1.24.5-alpine3.22
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY .. .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o api-server
EXPOSE 3000
CMD ["./api-server"]