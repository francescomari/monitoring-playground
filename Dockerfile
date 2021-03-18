FROM golang:1.16.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine:3.13.2
COPY  --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
