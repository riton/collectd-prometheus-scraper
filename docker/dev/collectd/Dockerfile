ARG COLLECTD_VERSION
FROM gitlab-registry.in2p3.fr/cc-in2p3/collectd-pkg/go-collectd:${COLLECTD_VERSION:-5.9.1} AS builder

ADD . /build/

WORKDIR /build/

RUN GOMODULES=on GOFLAGS="-mod=vendor" GOOS=linux go build -buildmode=c-shared -o prometheus_scraper.so

##
# Quick and dirty app
##
FROM gitlab-registry.in2p3.fr/cc-in2p3/containers/collectd:latest

USER root:root

COPY --from=builder /build/prometheus_scraper.so /usr/lib64/collectd/prometheus_scraper.so
COPY docker/dev/collectd/prometheus-LoadPlugin.conf "${COLLECTD_PLUGINS_CFG_DIR}/prometheus_scraper.conf"
COPY docker/dev/collectd/collectd-prometheus_scraper.yaml /etc/collectd-prometheus_scraper.yaml

RUN chown -R collectd:collectd /etc/collectd* && \
mkdir -p /var/run/collectd && \
chown -R collectd:collectd /var/run/collectd && \
chmod 0644 /etc/collectd-prometheus_scraper.yaml

USER collectd:collectd
