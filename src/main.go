package main

import (
	"dns-publisher/publisher"
	"flag"
	"net"
	"os"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var (
	configPathOpt = flag.String("configPath", "config.json", "Path to configuration file")
	logLevelOpt   = flag.String("logLevel", "INFO", "Set log level (NONE, ERROR, WARN, INFO, DEBUG)")
)

func main() {
	flag.Parse()

	loglevel, err := boshlog.Levelify(*logLevelOpt)
	if err != nil {
		loglevel = boshlog.LevelError
	}

	logger := boshlog.NewLogger(loglevel)

	config, err := NewConfigFromPath(*configPathOpt)
	if err != nil {
		logger.Error("main", "Loading config %s", err.Error())
		os.Exit(1)
	}
	logger.Info("main", "Configuration loaded")

	ticker := time.Tick(config.duration)

	publisher, err := publisher.NewPublisher(config.Publish)
	if err != nil {
		logger.Error("main", "Determining publisher - %s", err.Error())
		os.Exit(1)
	}

	data, err := publisher.Current()
	if err != nil {
		logger.Error("main", "Retrieving current configuration - %s", err.Error())
		os.Exit(1)
	}
	logger.Info("main", "Statup state includes %d entries.\n", len(data))

	for range ticker {
		// check and refresh
		logger.Info("main", "Updating from DNS")
		for query, _ := range config.DNS.ByQuery {
			ips, err := net.LookupIP(query)
			if err != nil {
				logger.Warn("main", "unable to lookup '%s': %s", query, err.Error())
				continue
			}
			logger.Debug("main", "found '%s' is %v", query, ips)
		}
	}
}
