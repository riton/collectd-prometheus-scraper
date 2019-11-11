package main

import (
	"fmt"

	"collectd.org/plugin"
	"github.com/ccin2p3/collectd-prometheus-plugin/logging"
)

const (
	pluginName = "prometheus_scraper"
)

var (
	gConfig   *config
	logger    logging.Logger
	underTest = false
)

func main() {} // unused

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

	if len(gConfig.ScrapeConfigs) == 0 {
		logger.Error("empty scrape_config")
		panic(fmt.Errorf("empty scrape_config"))
	}

	if gConfig.Debug {
		logger.SetDebug(true)
	}

	logger.Debugf("configuration: %+v", gConfig)

	//pscraper := scraper.NewPrometheusScraper("traefik")

	// register as many read callbacks as scrape targets we have
	for _, scrapeConfig := range gConfig.ScrapeConfigs {
		pscraper, err := newPrometheusScraper(scrapeConfig)
		if err != nil {
			logger.Errorf("initializing scraper %s, skipping", scrapeConfig.JobName)
			continue
		}

		readFunctionName := pluginName + "_" + scrapeConfig.JobName
		plugin.RegisterRead(readFunctionName, pscraper)
	}
}
