package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type teleinfoCollector struct {
	tic_isoucs *prometheus.Desc
	tic_index  *prometheus.Desc
	tic_papp   *prometheus.Desc
	tic_iinst  *prometheus.Desc
}

var teleinfoMetrics = TeleinfoMetrics{metric: make(map[string]string)}

func newTeleinfoCollector() *teleinfoCollector {
	return &teleinfoCollector{
		tic_isoucs: prometheus.NewDesc("tic_isoucs", "Inntensité souscrite", []string{"adco"}, nil),
		tic_index:  prometheus.NewDesc("tic_index", "Index", []string{"adco", "option", "color", "phase"}, nil),
		tic_papp:   prometheus.NewDesc("tic_papp", "Puissance apparente", []string{"adco"}, nil),
		tic_iinst:  prometheus.NewDesc("tic_iinst", "Intensite instantanee", []string{"adco"}, nil),
	}
}

func (collector *teleinfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.tic_isoucs
	ch <- collector.tic_index
	ch <- collector.tic_papp
	ch <- collector.tic_iinst
}

func (collector *teleinfoCollector) Collect(ch chan<- prometheus.Metric) {
	adco := teleinfoMetrics.Get("ADCO")
	if isoucs, err := strconv.ParseFloat(teleinfoMetrics.Get("ISOUCS"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_isoucs, prometheus.CounterValue, isoucs, adco)
	}
	if base, err := strconv.ParseFloat(teleinfoMetrics.Get("BASE"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, base, adco, "base", "none", "none")
	}
	if bbrhcjb, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHCJB"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhcjb, adco, "tempo", "blue", "hc")
	}
	if bbrhpjb, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHPJB"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhpjb, adco, "tempo", "blue", "hp")
	}
	if bbrhcjw, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHCJW"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhcjw, adco, "tempo", "white", "hc")
	}
	if bbrhpjw, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHPJW"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhpjw, adco, "tempo", "white", "hp")
	}
	if bbrhcjr, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHCJR"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhcjr, adco, "tempo", "red", "hc")
	}
	if bbrhpjr, err := strconv.ParseFloat(teleinfoMetrics.Get("BBRHPJR"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_index, prometheus.CounterValue, bbrhpjr, adco, "tempo", "red", "hp")
	}
	if papp, err := strconv.ParseFloat(teleinfoMetrics.Get("PAPP"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_papp, prometheus.GaugeValue, papp, adco)
	}
	if iinst, err := strconv.ParseFloat(teleinfoMetrics.Get("IINST"), 64); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.tic_iinst, prometheus.GaugeValue, iinst, adco)
	}
}

func main() {
	log.Printf("starting teleinfo_exporter")
	var (
		promPort = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)
	flag.Parse()

	// Execute the routine to get metrics from serial port
	log.Printf("starting go routine to grab information on serial port")
	go getSerialTeleinfo(&teleinfoMetrics)

	log.Printf("fetching initial data")
	time.Sleep(time.Duration(3) * time.Second)

	log.Printf("starting prometheus collector")
	teleinfoCollector := newTeleinfoCollector()

	log.Printf("adding prometheus collector to prometheus registry")
	reg := prometheus.NewRegistry()
	//reg.MustRegister(collectors.NewGoCollector()) // internal go process stats
	reg.MustRegister(teleinfoCollector)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting http server on %q", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start teleinfo_exporter: %s", err)
	}

}
