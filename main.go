package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var configFile string
var workerInterval int
var logLevel uint
func main() {

	flag.StringVar(&configFile,"c","conf/aa.yml","config file path")
	flag.IntVar(&workerInterval,"i",120,"data scrape interval")
	flag.UintVar(&logLevel,"l",3,"loglevel 1:fatal 2:error 3:warn 4:info 5:debug 6:trace")
	flag.Parse()

	log.SetLevel(log.Level(logLevel))
	//Create a new instance of the foocollector and
	//register it with the prometheus client.
	foo := newApiCollector()
	prometheus.MustRegister(foo)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Warn("Beginning to serve on port :10880")
	log.Fatal(http.ListenAndServe(":10880", nil))
}