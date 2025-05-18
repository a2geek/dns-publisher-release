package web

import (
	"dns-publisher/config"
	"fmt"
	"log"
	"net/http"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var configuration config.Config

func Server(config config.Config, logger boshlog.Logger) {
	configuration = config
	http.HandleFunc("/api/logs", getLogs)
	http.HandleFunc("/api/config", getConfig)
	http.HandleFunc("/api/tags", getTags)
	// TODO be smart about port
	logger.Info("web-server", "Starting HTTP server on port %d", config.Web.HTTP)
	addr := fmt.Sprintf(":%d", config.Web.HTTP)
	log.Fatal(http.ListenAndServe(addr, nil))
}
