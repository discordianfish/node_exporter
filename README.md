# Supervisor exporter [![Build Status](https://travis-ci.org/prometheus/supervisor_exporter.svg)][travis]

[![CircleCI](https://circleci.com/gh/prometheus/supervisor_exporter/tree/master.svg?style=shield)][circleci]
[![Docker Repository on Quay](https://quay.io/repository/prometheus/supervisor-exporter/status)][quay]
[![Docker Pulls](https://img.shields.io/docker/pulls/prom/supervisor-exporter.svg?maxAge=604800)][hub]

Prometheus exporter for supervisor metrics, written in Go with pluggable metric
collectors.

## Collectors

There is varying support for supervisor on each operating system. The tables
below list all existing collectors and the supported systems.

Which collectors are used is controlled by the `--collectors.enabled` flag.

### Enabled by default

Name     | Description | OS
---------|-------------|----
logind | Exposes session counts from [logind](http://www.freedesktop.org/wiki/Software/systemd/logind/). | Linux
runit | Exposes service status from [runit](http://smarden.org/runit/). | _any_
supervisord | Exposes service status from [supervisord](http://supervisord.org/). | _any_
systemd | Exposes service and system status from [systemd](http://www.freedesktop.org/wiki/Software/systemd/). | Linux

## Building and running

    make
    ./supervisor_exporter <flags>

## Running tests

    make test

[travis]: https://travis-ci.org/prometheus/supervisor_exporter
[hub]: https://hub.docker.com/r/prom/supervisor-exporter/
[circleci]: https://circleci.com/gh/prometheus/supervisor_exporter
[quay]: https://quay.io/repository/prometheus/supervisor-exporter
