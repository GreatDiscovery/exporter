package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
)
//简单的exporter，自定义两个metric，采集CPU温度和磁盘损坏次数，拉取的数据由自己给出
//采集器的结构体
type Collector struct{
	Zone string
	cpuTemperature *prometheus.Desc
	hdFailure *prometheus.Desc
}

func (c *Collector) CollectorImpl()(cpuByHost map[string]int, hdByHost map[string]int,){
	cpuByHost = map[string]int{
		"localhost1":int(rand.Int31n(100)),
		"localhost2":int(rand.Int31n(100)),
	}
	hdByHost = map[string]int{
		"localhost1":int(rand.Int31n(100)),
		"localhost2":int(rand.Int31n(100)),
	}
	return
}
//实现collector接口
func (c *Collector) Describe(ch chan <- *prometheus.Desc){
	ch <- c.cpuTemperature
	ch <- c.hdFailure
}


func(c *Collector) Collect(ch chan <- prometheus.Metric) {
	cpuByHost,hdByHost:=c.CollectorImpl()
	for host, cpuCount:=range cpuByHost{
		ch <- prometheus.MustNewConstMetric(
			c.cpuTemperature,
			prometheus.CounterValue,
			float64(cpuCount),
			host,)
	}
	for host,hdCount:= range hdByHost{
		ch <- prometheus.MustNewConstMetric(
			c.hdFailure,
			prometheus.GaugeValue,
			float64(hdCount),
			host,
		)
	}
}
//重新定义采集器
func NewCollector(zone string) *Collector{
	return &Collector{
		Zone:zone,
		cpuTemperature:prometheus.NewDesc(
			"cpuTemperature_total",
			"Number of cpu's temperature.",
			[]string{"host"},
			prometheus.Labels{"zone":zone},
		),
		hdFailure:prometheus.NewDesc(
			"hdFailure_total",
			"Number of hd's failure",
			[]string{"host"},
			prometheus.Labels{"zone":zone},
		),
	}
}
func main(){
	workerDB := NewCollector("db")
	workerCA := NewCollector("ca")
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workerDB)
	reg.MustRegister(workerCA)
	gatherers:=prometheus.Gatherers{
		prometheus.DefaultGatherer,
		reg,
	}
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
		})
	http.HandleFunc("/metrics",func(w http.ResponseWriter,r *http.Request){
		h.ServeHTTP(w,r)
	})
	http.ListenAndServe(":8081", nil)

}
