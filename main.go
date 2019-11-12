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

	"collectd.org/plugin"
	"gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/logging"
)

const (
	pluginName                = "prometheus_scraper"
	collectdReadCallbackGroup = "prometheus_scraper"
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
		plugin.RegisterComplexRead(readFunctionName, pscraper, plugin.ComplexReadConfig{
			Group: collectdReadCallbackGroup,
		})
	}
}
