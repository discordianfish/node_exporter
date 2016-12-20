FROM        quay.io/prometheus/busybox:latest
MAINTAINER  The Prometheus Authors <prometheus-developers@googlegroups.com>

COPY supervisor_exporter /bin/supervisor_exporter

EXPOSE      9228
ENTRYPOINT  [ "/bin/supervisor_exporter" ]
