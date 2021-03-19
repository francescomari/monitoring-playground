# Monitoring Playground

This repository is a playground to experiment with a monitoring stack. Using
Docker Compose, you can deploy:

- [Prometheus](http://localhost:9090)
- [Alertmanager](http://localhost:9093)
- [Grafana](http://localhost:3000)
- A [sample application](http://localhost:8080/metrics)

The [configuration for Prometheus](prometheus) is mounted as a volume and
consumed by the corresponding service running the official Prometheus Docker
image. Prometheus has been configured to scrape all the services run in Docker
Compose, including itself.

The [configuration for Grafana](grafana) is mounted as a volume and consumed by
the corresponding service running the official Grafana Docker image.

The [sample application](main.go) is automatically package in a Docker image and
run as a service in Docker Compose. The sample application has no configuration,
and generates random request metrics.

## Usage

Start Docker Compose.

    make up

Stop the containers started by Docker Compose and remove all the Docker
resources previously created (containers, networks, etc.).

    make down

Build the sample application. This is only necessary if you change the sample
application and want to re-package it before starting Docker Compose.

    make build

## Sample application

The sample application pretends to continuously receive a certain amount of
requests and exposes two metrics related to these requests:

- `app_request_duration_seconds` - histogram - The duration of the requests, in
  seconds.
- `app_request_errors_count` - counter - The number of requests resulting in an
  error.

By default, the sample application will randomly observe a request duration
between 0s and 10s and will simulate 0.1 errors/request. Moreover, the initial
request rate is 1 request/s. These limits can be changed at runtime using
`curl`.

    curl -X PATCH localhost:8080/limits -F maxDuration=20.0 -F errorsRatio=0.5 -F requestRate=1.5

The request above will instruct the sample application to observe a request
duration between 0s and 20s, and to simulate 0.5 errors/request. The service is
also instructed to pretend it's receiving requests with a rate of 1.5
requests/s.

`maxDuration`, `errorsRatio`, and `requestRate` must be floating point numbers.
`maxDuration` must be a number greater than or equal to zero. `errorsRatio` must
be a number between 0 and 1, inclusive. `requestRate` must be strictly greater
than zero. You can specify every parameter at once, none, or any combination of
them. If specified, the parameters must be valid.

If the limits are valid, they are immediately applied. In this case, the sample
application returns a 200 response and prints a log message, which will be
visible in the output of Docker Compose. If the request is not valid or any
error occurs while processing the request, the sample application returns a
4xx or 5xx response and prints an error message in its logs.
