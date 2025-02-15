/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	etcdRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "etcd_request_duration_seconds",
			Help: "Etcd request latency in seconds for each operation and object type.",
		},
		[]string{"operation", "type"},
	)
	objectCounts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "etcd_object_counts",
			Help: "Number of stored objects at the time of last check split by kind.",
		},
		[]string{"resource"},
	)

	deprecatedEtcdRequestLatenciesSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "etcd_request_latencies_summary",
			Help: "(Deprecated) Etcd request latency summary in microseconds for each operation and object type.",
		},
		[]string{"operation", "type"},
	)
)

var registerMetrics sync.Once

// Register all metrics.
func Register() {
	// Register the metrics.
	registerMetrics.Do(func() {
		prometheus.MustRegister(etcdRequestLatency)
		prometheus.MustRegister(objectCounts)

		// TODO(danielqsj): Remove the following metrics, they are deprecated
		prometheus.MustRegister(deprecatedEtcdRequestLatenciesSummary)
	})
}

func UpdateObjectCount(resourcePrefix string, count int64) {
	objectCounts.WithLabelValues(resourcePrefix).Set(float64(count))
}

func RecordEtcdRequestLatency(verb, resource string, startTime time.Time) {
	etcdRequestLatency.WithLabelValues(verb, resource).Observe(sinceInSeconds(startTime))
	deprecatedEtcdRequestLatenciesSummary.WithLabelValues(verb, resource).Observe(sinceInMicroseconds(startTime))
}

func Reset() {
	etcdRequestLatency.Reset()

	deprecatedEtcdRequestLatenciesSummary.Reset()
}

// sinceInMicroseconds gets the time since the specified start in microseconds.
func sinceInMicroseconds(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds() / time.Microsecond.Nanoseconds())
}

// sinceInSeconds gets the time since the specified start in seconds.
func sinceInSeconds(start time.Time) float64 {
	return time.Since(start).Seconds()
}
