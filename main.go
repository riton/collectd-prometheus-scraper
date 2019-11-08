package main

import (
	"collectd.org/plugin"
	"gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/logging"
	"gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/scraper"
)

const (
	pluginName = "prometheus_scraper"
)

var (
	gConfig   *config
	logger    logging.Logger
	underTest = false
)

func main() {
	scraper := scraper.NewPrometheusScraper("traefik")
	scraper.Read()
}

func init() {
	if !underTest {
		_init()
	}
}

func _init() {
	gConfig = newConfig()
	logger = logging.NewCollectdLogger(gConfig.Debug)

	pscraper := scraper.NewPrometheusScraper("traefik")

	plugin.RegisterRead(pluginName, pscraper)
}
