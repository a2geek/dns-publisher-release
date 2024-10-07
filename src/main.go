package main

import (
	"dns-publisher/processors"
	"dns-publisher/publishers"
	"flag"
	"os"

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

	publisher, err := publishers.NewPublisher(config.Publisher, logger)
	if err != nil {
		logger.Error("main", "Determining publisher - %s", err.Error())
		os.Exit(1)
	}

	if config.BoshDns != nil {
		processor, err := processors.NewBoshDnsProcessor(*config.BoshDns, publisher, logger)
		if err != nil {
			logger.Error("main", "Unable to create BOSH DNS processor: %s", err.Error())
			os.Exit(1)
		}
		go processor.Run()
	}

	// wait forever
	select {}
}
