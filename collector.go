package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type apiCollector struct {
	httpServiceUp *prometheus.Desc
	httpServiceTime *prometheus.Desc
	tcpServiceUp *prometheus.Desc
	workerBean workerBean
}

func newApiCollector() *apiCollector {
	apiCollectorIns := apiCollector{
		httpServiceUp: prometheus.NewDesc("http_service_up",
			"service is up",
			[]string{"name","url"}, nil, ),
		httpServiceTime: prometheus.NewDesc("http_service_time",
			"response time",[]string{"name","url"},nil),
		tcpServiceUp: prometheus.NewDesc("tcp_service_up",
				"service is up",
				[]string{"name","host","port"},nil),
	}
	wb,err := initWorker()
	if err != nil{
		panic("init config failed")
	}else{
		apiCollectorIns.workerBean = wb
	}
	return &apiCollectorIns
}

func (collector *apiCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.httpServiceUp
	ch <- collector.httpServiceTime
	ch <- collector.tcpServiceUp
}

func (collector *apiCollector) Collect(ch chan<- prometheus.Metric) {

	wb := collector.workerBean
	hw := wb.Http
	for _,v := range hw{
		res,err := cache.Value(v.Name)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(collector.httpServiceUp, prometheus.GaugeValue,0,v.Name,v.Url)
			log.Debug("get http worker metirc err",v,err)
		}else{
			r := res.Data().(*Result)
			ch <- prometheus.MustNewConstMetric(collector.httpServiceUp, prometheus.GaugeValue,float64(r.Success),v.Name,v.Url)
			ch <- prometheus.MustNewConstMetric(collector.httpServiceTime, prometheus.GaugeValue,r.ResponseTime,v.Name,v.Url)
		}
	}
	tcp := wb.Tcp
	for _,v := range tcp {
		res,err := cache.Value(v.Name)
		if err != nil {
			log.Debug("get http worker metirc err",v,err)
			ch <- prometheus.MustNewConstMetric(collector.tcpServiceUp,prometheus.GaugeValue,0,v.Name,v.Host,strconv.Itoa(v.Port))
		}else{
			b := res.Data().(*Result)
			ch <- prometheus.MustNewConstMetric(collector.tcpServiceUp,prometheus.GaugeValue,float64(b.Success),v.Name,v.Host,strconv.Itoa(v.Port))
		}
	}



}
