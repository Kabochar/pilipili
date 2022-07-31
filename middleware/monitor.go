package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

/*
Gin监控第三方库：https://github.com/penglongli/gin-metrics

这里搭配 Prometheus，从官网下载并安装一个具体版本即可，具体操作流程跳过。

默认是不能直接监控直接服务数据，需要到 Prometheus 安装目录下找到 prometheus.yml，找到 static_configs 选项，加入的具体的监控地址
static_configs:
  - targets: ["localhost:9090"]
  - targets: ["localhost:9000"] // 比如我的是 localhost:9000/metrics

Prometheus原生的报表不够炫酷，所以这里接入 Grafana 来展示具体的数据。Grafana 配置具体的 Dashboard，这里第三方库提供了一个案例json
https://github.com/penglongli/gin-metrics/blob/master/grafana/grafana.json

DashBoard配置方式
选择 New Dashboard，import json，将上面链接的具体json内容复制进去即可。

注：
原来的报表没有将请求次数进行排序展示，这里作为监控并不太友好，可以留个todo自行补充
Prometheus官网下载地址：https://prometheus.io/download/
Grafana官网下载地址：https://grafana.com/grafana/download?edition=oss
第三方库Grafana需要的插件：https://grafana.com/grafana/plugins/grafana-piechart-panel/?tab=installation
Grafana插件安装教程：https://grafana.com/docs/grafana/next/cli/#plugins-commands
匿名访问Grafana：https://blog.csdn.net/qq_22227087/article/details/86993324
*/

// 接入Prometheus监控
func GinMetrics(engine *gin.Engine) {
	// get global Monitor object
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(engine)
}
