package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	taskQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transform_task_queue_size",
		Help: "Size of the task queue",
	})
)

var retryTimes float64 = 0
var success = false
var api = os.Getenv("WEB_API")
var duration float64 = 0

func init() {
	prometheus.MustRegister(taskQueueSize)
}

type ResponBody struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

func updateTaskQueueSize() {
	if api == "" {
		api = "http://localhost:7001/user/video/getQeuenLen?queueName=transform_4.0"
	}
	for {
		if !success {
			// 最长5分钟
			if retryTimes > 10 {
				retryTimes = 0
			} else {
				duration = math.Pow(2, retryTimes)
				retryTimes++
			}
			time.Sleep(time.Second * time.Duration(duration))
		}
		// 发起 HTTP 请求，获取任务数量
		resp, err := http.Get(api)

		if err != nil {
			fmt.Printf("Error occurred while fetching task count: %v\n", err)
			continue
		}
		// defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch task count. Status code: %v\n", resp.StatusCode)
			continue
		}

		// 读取响应体中的任务数量
		var result ResponBody
		err = json.NewDecoder(resp.Body).Decode(&result)

		if err != nil {
			fmt.Printf("Failed to parse task count: %v\n", err)
			continue
		}

		fmt.Printf("Update task count to: %v\n", result.Data)

		// 更新指标的值
		taskQueueSize.Set(float64(result.Data))
		success = true
		time.Sleep(5 * time.Second) // 每隔一分钟更新一次
	}
}

func main() {
	go updateTaskQueueSize()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8000", nil)
}
