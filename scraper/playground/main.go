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
