package main

import (
	"fmt"

	"collectd.org/api"
	"github.com/pkg/errors"

	"gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/scraper"
)

func newPrometheusScraper(cfg scraperConfig) (*scraper.PrometheusScraper, error) {

	additionalMeta := make(api.Metadata)
	for metaName := range cfg.Labels {
		additionalMeta.Set(metaName, cfg.Labels[metaName])
	}

	targetURL := fmt.Sprintf("%s://%s:%d%s", cfg.Scheme,
		cfg.Target, cfg.Port, cfg.MetricsPath)

	scraper := scraper.NewPrometheusScraper(
		cfg.Collectd.PluginName,
		cfg.Collectd.MetadataPrefix,
		targetURL,
		cfg.ScrapeTimeout,
		cfg.BasicAuth,
		cfg.Collectd.TypeInstanceOnlyHashedMeta,
		cfg.Collectd.HashLabelFunctionHashSize,
		additionalMeta,
	)

	if err := scraper.Initialize(); err != nil {
		return scraper, errors.Wrap(err, "initializing scraper")
	}

	return scraper, nil
}
