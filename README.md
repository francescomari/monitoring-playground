# Monitoring Playground

This repository is a playground to experiment with a monitoring stack. Using
Docker Compose, you can deploy:

- [Prometheus](http://localhost:9090)
- [Alertmanager](http://localhost:9093)
- [Grafana](http://localhost:3000)
- [Metrics Generator](http://localhost:8080/metrics)

The [configuration for Prometheus](prometheus) is mounted as a volume and
consumed by the corresponding service running the official Prometheus Docker
image. Prometheus has been configured to scrape all the services run in Docker
Compose, including itself.

The [configuration for Grafana](grafana) is mounted as a volume and consumed by
the corresponding service running the official Grafana Docker image.

[Metrics Generator](https://github.com/francescomari/metrics-generator) is a
service that pretends it's receiving requests and exposes metrics about those
request. It is included in the playground in order to have a constant source of
pseudo-random data to play with.

## Usage

Start Docker Compose.

    make up

Stop the containers started by Docker Compose and remove all the Docker
resources previously created (containers, networks, etc.).

    make down
