package main

import (
	"fmt"
	"strings"

	pexpfmt "github.com/prometheus/common/expfmt"
)

var promData = `
# HELP traefik_service_request_duration_seconds How long it took to process the request on a service, partitioned by status code, protocol, and method.
# TYPE traefik_service_request_duration_seconds histogram
traefik_service_request_duration_seconds_bucket{code="200",method="GET",protocol="http",service="whoami-traefik@docker",le="0.1"} 100428
traefik_service_request_duration_seconds_bucket{code="200",method="GET",protocol="http",service="whoami-traefik@docker",le="0.3"} 100428
traefik_service_request_duration_seconds_bucket{code="200",method="GET",protocol="http",service="whoami-traefik@docker",le="1.2"} 100428
traefik_service_request_duration_seconds_bucket{code="200",method="GET",protocol="http",service="whoami-traefik@docker",le="5"} 100428
traefik_service_request_duration_seconds_bucket{code="200",method="GET",protocol="http",service="whoami-traefik@docker",le="+Inf"} 100428
traefik_service_request_duration_seconds_sum{code="200",method="GET",protocol="http",service="whoami-traefik@docker"} 148.226708722
traefik_service_request_duration_seconds_count{code="200",method="GET",protocol="http",service="whoami-traefik@docker"} 100428
traefik_service_request_duration_seconds_bucket{code="499",method="GET",protocol="http",service="whoami-traefik@docker",le="0.1"} 1
traefik_service_request_duration_seconds_bucket{code="499",method="GET",protocol="http",service="whoami-traefik@docker",le="0.3"} 1
traefik_service_request_duration_seconds_bucket{code="499",method="GET",protocol="http",service="whoami-traefik@docker",le="1.2"} 1
traefik_service_request_duration_seconds_bucket{code="499",method="GET",protocol="http",service="whoami-traefik@docker",le="5"} 1
traefik_service_request_duration_seconds_bucket{code="499",method="GET",protocol="http",service="whoami-traefik@docker",le="+Inf"} 1
traefik_service_request_duration_seconds_sum{code="499",method="GET",protocol="http",service="whoami-traefik@docker"} 0.000444565
traefik_service_request_duration_seconds_count{code="499",method="GET",protocol="http",service="whoami-traefik@docker"} 1
`

func main() {

	promDataIOReader := strings.NewReader(promData)

	tparser := pexpfmt.TextParser{}

	metrics, err := tparser.TextToMetricFamilies(promDataIOReader)
	if err != nil {
		panic(err)
	}

	for _, metric := range metrics {
		fmt.Printf("%+v\n", metric)
	}

	fmt.Printf("\n-- Metrics --\n\n")

	for _, mf := range metrics {
		for idx, metric := range mf.Metric {
			fmt.Printf("[METRIC_%d] %+v\n\n", idx, metric)
		}
	}
}
