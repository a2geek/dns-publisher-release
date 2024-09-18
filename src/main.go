package main

import (
	"dns-publisher/filter"
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

	logger boshlog.Logger
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

	filters := config.Filters.GetFilters()
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
	applyFilters(data, filters)
	logger.Info("main", "Startup state includes %d entries.\n", len(data))

	for range ticker {
		// check and refresh
		logger.Info("main", "Updating from DNS")
		for query := range config.DNS.ByQuery {
			ips, err := net.LookupIP(query)
			if err != nil {
				logger.Warn("main", "unable to lookup '%s': %s", query, err.Error())
				continue
			}
			logger.Debug("main", "found '%s' is %v", query, ips)
		}
	}
}

func applyFilters(data map[string][]net.IP, filters []filter.IPFilter) {
	for host, ipaddrs := range data {
		var newaddrs []net.IP
		// rebuild the IP addr list
		for _, ipaddr := range ipaddrs {
			good := false
			for _, filter := range filters {
				good = good || filter(ipaddr)
			}
			if good {
				newaddrs = append(newaddrs, ipaddr)
			}
		}
		// handle map, depending on how many passed the filter
		if len(newaddrs) == 0 {
			delete(data, host)
			logger.Debug("Removing host '%s' from list as it did not pass filters", host)
		} else {
			data[host] = newaddrs
		}
	}
}
