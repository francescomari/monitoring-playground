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

### CLI

The sample application accepts flags to initialize the rate of the simulated
request, the maximum request duration and the percentage of requests that will
result in an error. Use the `-help` flag to see the command's help.

### API

  GET /-/health

Always return a 200 response.

  PUT /-/config/max-duration

Set the maximum duration of the simulated requests to the value passed in the
body of the request. The value must be a  positive integer.

  PUT /-/config/errors-percentage

Set the percentage of the simulated requests that will result in an error to the
value passed in the body of the request. It must be an integer between 0 and
100.

  PUT /-/config/request-rate

Set the rate of the simulated requests to the value passed in the body of the
request. It must be a strictly positive integer.
