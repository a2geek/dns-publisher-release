package main

import (
	"dns-publisher/publishers"
	"dns-publisher/sources"
	"flag"
	"net"
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var (
	configPathOpt = flag.String("configPath", "config.json", "Path to configuration file")
	logLevelOpt   = flag.String("logLevel", "INFO", "Set log level (NONE, ERROR, WARN, INFO, DEBUG)")

	logger    boshlog.Logger
	publisher publishers.Publisher
)

func main() {
	flag.Parse()

	loglevel, err := boshlog.Levelify(*logLevelOpt)
	if err != nil {
		loglevel = boshlog.LevelError
	}

	logger = boshlog.NewLogger(loglevel)

	config, err := NewConfigFromPath(*configPathOpt)
	if err != nil {
		logger.Error("main", "Loading config %s", err.Error())
		os.Exit(1)
	}
	logger.Info("main", "Configuration loaded")

	publisher, err = publishers.NewPublisher(config.Publisher, logger)
	if err != nil {
		logger.Error("main", "Determining publisher - %s", err.Error())
		os.Exit(1)
	}

	state, err := publisher.Current()
	if err != nil {
		logger.Error("main", "Retrieving current configuration - %s", err.Error())
		os.Exit(1)
	}
	logger.Info("main", "Startup state includes %d entries: %v", len(state), hostKeysAsString(state))

	source, err := sources.NewSource(config.Source)
	if err != nil {
		logger.Error("main", "Configuring source: %s", err.Error())
		os.Exit(1)
	}

	for range source.Start() {
		// check and refresh
		logger.Info("main", "Updating from DNS")

		state, err = publisher.Current()
		if err != nil {
			logger.Error("main", "Retrieving current configuration - %s", err.Error())
			os.Exit(1)
		}
		logger.Info("main", "Current state includes %d entries: %v\n", len(state), hostKeysAsString(state))
		changes := false
		for query, hosts := range config.Source.ByQuery {
			ips, err := source.Lookup(query)
			if err != nil {
				logger.Warn("main", "unable to lookup '%s': %s", query, err.Error())
				continue
			}
			logger.Debug("main", "found '%s' for %v = %v", query, hosts, ips)
			for _, host := range hosts {
				change, err := adjustState(state, host, ips)
				if err != nil {
					logger.Error("main", "error adjusting state for '%s': %v", host, err)
				}
				changes = changes || change
			}
		}
		if changes {
			err = publisher.Commit()
			if err != nil {
				logger.Error("main", "unable to commit changes: %s", err.Error())
			}
		}
	}
}

func adjustState(state map[string][]net.IP, host string, newIPs []net.IP) (bool, error) {
	currentIPs, ok := state[host]
	if ok && !sameContents(currentIPs, newIPs) {
		err := publisher.Delete(host)
		if err != nil {
			return false, err
		}
		delete(state, host)
	}
	if ok {
		return false, nil
	}
	err := publisher.Add(host, newIPs)
	if err != nil {
		return true, err
	}
	state[host] = newIPs
	return true, nil
}
