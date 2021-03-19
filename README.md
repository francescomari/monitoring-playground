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
