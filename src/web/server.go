package web

import (
	"dns-publisher/config"
	"fmt"
	"log"
	"net/http"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func NewWebServer(config config.Config, logger boshlog.Logger, logCache *LogCache) WebServer {
	return WebServer{
		config:   config,
		logger:   logger,
		logCache: logCache,
	}
}

type WebServer struct {
	config   config.Config
	logger   boshlog.Logger
	logCache *LogCache
}

func (s WebServer) Serve() {
	http.HandleFunc("/api/logs", s.getLogs)
	http.HandleFunc("/api/config", s.getConfig)
	http.HandleFunc("/api/tags", s.getTags)
	// TODO be smart about port
	s.logger.Info("web-server", "Starting HTTP server on port %d", s.config.Web.HTTP)
	addr := fmt.Sprintf(":%d", s.config.Web.HTTP)
	log.Fatal(http.ListenAndServe(addr, webLogger{s.logger}))
}

type webLogger struct {
	logger boshlog.Logger
}

func (l webLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	crw := NewCapturingResponseWriter(w)
	http.DefaultServeMux.ServeHTTP(crw, r)

	username, _, ok := r.BasicAuth()
	if !ok {
		username = "-"
	}

	remoteIP := strings.Split(r.RemoteAddr, ":")

	l.logger.Info("web-server", "%s %s \"%s %s %s\" %d %d",
		remoteIP[0], username, r.Method, r.URL.Path, r.Proto,
		crw.statusCode, crw.contentLength)
}

// See https://stackoverflow.com/questions/53272536/how-do-i-get-response-statuscode-in-golang-middleware
type capturingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func NewCapturingResponseWriter(w http.ResponseWriter) *capturingResponseWriter {
	return &capturingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *capturingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *capturingResponseWriter) Write(b []byte) (int, error) {
	lrw.contentLength += len(b)
	return lrw.ResponseWriter.Write(b)
}
