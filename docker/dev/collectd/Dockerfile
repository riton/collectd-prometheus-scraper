ARG COLLECTD_VERSION
FROM gitlab-registry.in2p3.fr/cc-in2p3/collectd-pkg/go-collectd:${COLLECTD_VERSION:-5.9.1} AS builder

ADD . /build/

WORKDIR /build/

RUN GOMODULES=on GOFLAGS="-mod=vendor" GOOS=linux go build -buildmode=c-shared -o prometheus_scraper.so

##
# Quick and dirty app
##
FROM gitlab-registry.in2p3.fr/rferrand/collectd/feature-meta

USER root:root

COPY --from=builder /build/prometheus_scraper.so /src/install/lib/collectd/prometheus_scraper.so
COPY docker/dev/collectd/minimal_config.conf /etc/collectd.d/minimal_config.conf

RUN chown -R collectd:collectd /etc/collectd* && \
mkdir -p /var/run/collectd && \
chown -R collectd:collectd /var/run/collectd

USER collectd:collectd
