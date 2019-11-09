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
	scraper := scraper.NewPrometheusScraper("coredns")
	//scraper := scraper.NewPrometheusScraper("traefik")
	scraper.Read()
}

func init() {
	if !underTest {
		_init()
	}
}

func _init() {

	logger = logging.NewCollectdLogger(pluginName + ": " /* logPrefix */)

	gConfig, err := newConfig()
	if err != nil {
		logger.Errorf("reading configuration file %s", err)
		panic(err)
	}

	if gConfig.Debug {
		logger.SetDebug(true)
	}

	logger.Debugf("configuration: %+v", gConfig)

	pscraper := scraper.NewPrometheusScraper("coredns")
	//pscraper := scraper.NewPrometheusScraper("traefik")

	plugin.RegisterRead(pluginName, pscraper)
}
