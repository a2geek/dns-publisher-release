package main

import (
	"dns-publisher/publishers"
	"flag"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var (
	configPathOpt = flag.String("configPath", "config.json", "Path to configuration file")
	logLevelOpt   = flag.String("logLevel", "INFO", "Set log level (NONE, ERROR, WARN, INFO, DEBUG)")
	dryRunOpt     = flag.Bool("dryRun", false, "Disallow any change operations")

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

	ticker := time.Tick(config.duration)

	publisher, err = publishers.NewPublisher(config.Publish, logger, *dryRunOpt)
	if err != nil {
		logger.Error("main", "Determining publisher - %s", err.Error())
		os.Exit(1)
	}

	state, err := publisher.Current()
	if err != nil {
		logger.Error("main", "Retrieving current configuration - %s", err.Error())
		os.Exit(1)
	}
	logger.Info("main", "Startup state includes %d entries: %v\n", len(state), hostKeysAsString(state))

	for range ticker {
		// check and refresh
		logger.Info("main", "Updating from DNS")

		state, err = publisher.Current()
		if err != nil {
			logger.Error("main", "Retrieving current configuration - %s", err.Error())
			os.Exit(1)
		}
		logger.Info("main", "Current state includes %d entries: %v\n", len(state), hostKeysAsString(state))

		for query, hosts := range config.DNS.ByQuery {
			ips, err := net.LookupIP(query)
			if err != nil {
				logger.Warn("main", "unable to lookup '%s': %s", query, err.Error())
				// TODO likely want a flag to either delete or leave as-is if host exists in state
				continue
			}
			logger.Debug("main", "found '%s' for %v = %v", query, hosts, ips)
			for _, host := range hosts {
				err = adjustState(state, host, ips)
				if err != nil {
					logger.Error("main", "error adjusting state for '%s': %v", host, err)
				}
			}
		}
	}
}

func adjustState(state map[string][]net.IP, host string, newIPs []net.IP) error {
	currentIPs, ok := state[host]
	if ok && !reflect.DeepEqual(currentIPs, newIPs) {
		err := publisher.Delete(host)
		if err != nil {
			return err
		}
		delete(state, host)
	}
	err := publisher.Add(host, newIPs)
	if err != nil {
		return err
	}
	state[host] = newIPs
	return err
}

func hostKeysAsString(state map[string][]net.IP) string {
	var builder strings.Builder
	for key := range state {
		if builder.Len() > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(key)
	}
	return builder.String()
}
