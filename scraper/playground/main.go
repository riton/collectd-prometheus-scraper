// Copyright (c) IN2P3 Computing Centre, IN2P3, CNRS
// 
// Author(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2019
// 
// This software is governed by the CeCILL-C license under French law and
// abiding by the rules of distribution of free software.  You can  use, 
// modify and/ or redistribute the software under the terms of the CeCILL-C
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info". 
// 
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability. 
// 
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or 
// data to be ensured and,  more generally, to use and operate it in the 
// same conditions as regards security. 
// 
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL-C license and that you accept its terms.

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/ccin2p3/collectd-prometheus-plugin/transport"
)

// All this configuration stuff should be integrated in
// collectd complex configuration callback as soon as the
// Go bindings are available

const (
	configEnvVarPrefix     = "COLLECTD_PROM"
	configBaseDir          = "."
	configFilename         = "collectd-prometheus_scraper"
	defaultScrapeTimeout   = 5 * time.Second
	defaultMetricsPath     = "/metrics"
	defaultHonorTimestamps = false
)

// scrapeConfig is heavily inspired by
// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config
type scraperConfig struct {
	JobName         string                   `mapstructure:"job_name"`
	Target          string                   `mapstructure:"target"`
	HonorTimestamps bool                     `mapstructure:"honor_timestamps"`
	ScrapeTimeout   time.Duration            `mapstructure:"scrape_timeout"`
	MetricsPath     string                   `mapstructure:"metrics_path"`
	Labels          map[string]interface{}   `mapstructure:"labels"`
	BasicAuth       transport.HTTPBasicCreds `mapstructure:"basic_auth"`
}

type config struct {
	Debug         bool            `mapstructure:"debug"`
	ScrapeTimeout time.Duration   `mapstructure:"scrape_timeout"`
	ScrapeConfigs []scraperConfig `mapstructure:"scrape_configs"`
}

func main() {
	cfg, err := newConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}

func newConfig() (config, error) {
	var cfg config

	vcfg, err := newViperConfig()
	if err != nil {
		return cfg, errors.Wrap(err, "initializing viper")
	}

	cfg.Debug = vcfg.GetBool("debug")
	cfg.ScrapeTimeout = vcfg.GetDuration("scrape_timeout")

	if !vcfg.IsSet("scrape_configs") {
		return cfg, errors.New("missing scrape_configs section")
	}

	for scrapeKey := range vcfg.GetStringMap("scrape_configs") {

		sub := vcfg.Sub("scrape_configs." + scrapeKey)
		sub.SetDefault("scrape_timeout", cfg.ScrapeTimeout)
		sub.SetDefault("metrics_path", defaultMetricsPath)
		sub.SetDefault("honor_timestamps", defaultHonorTimestamps)

		var sconfig scraperConfig

		sconfig.JobName = scrapeKey
		if err := sub.Unmarshal(&sconfig); err != nil {
			return cfg, errors.Wrap(err, fmt.Sprintf("unmarsharling configuration for %s",
				scrapeKey))
		}

		cfg.ScrapeConfigs = append(cfg.ScrapeConfigs, sconfig)
	}

	return cfg, nil
}

func newViperConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(configBaseDir)
	v.SetConfigName(configFilename)

	v.SetEnvPrefix(configEnvVarPrefix)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	v.AutomaticEnv()

	v.SetDefault("debug", false)
	v.SetDefault("scrape_timeout", defaultScrapeTimeout)

	err := v.ReadInConfig()

	return v, err
}
