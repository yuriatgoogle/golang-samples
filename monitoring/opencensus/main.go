// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// [START monitoring_sli_metrics_opencensus_setup]
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"golang.org/x/exp/rand"
)

// [END monitoring_sli_metrics_opencensus_setup]
// [START monitoring_sli_metrics_opencensus_measure]
// set up metrics
var (
	requestCount       = stats.Int64("oc_request_count", "total request count", "requests")
	failedRequestCount = stats.Int64("oc_failed_request_count", "count of failed requests", "requests")
	responseLatency    = stats.Float64("oc_latency_distribution", "distribution of response latencies", "s")
)

// [END monitoring_sli_metrics_opencensus_measure]

// set up views
// [START monitoring_sli_metrics_opencensus_view]
var (
	requestCountView = &view.View{
		Name:        "oc_request_count",
		Measure:     requestCount,
		Description: "total request count",
		Aggregation: view.Count(),
	}
	failedRequestCountView = &view.View{
		Name:        "oc_failed_request_count",
		Measure:     failedRequestCount,
		Description: "count of failed requests",
		Aggregation: view.Count(),
	}
	responseLatencyView = &view.View{
		Name:        "oc_response_latency",
		Measure:     responseLatency,
		Description: "The distribution of the latencies",
		// bucket definitions must be explicit
		Aggregation: view.Distribution(0, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000),
	}
)

// [END monitoring_sli_metrics_opencensus_view]

func main() {

	// set up project ID
	projectID := flag.String("project_id", "", "Cloud Project ID")
	flag.Parse()

	// register the views
	if err := view.Register(requestCountView, failedRequestCountView, responseLatencyView); err != nil {
		log.Fatalf("Failed to register the views: %v", err)
	}
	// [START monitoring_sli_metrics_opencensus_exporter]
	// set up Cloud Monitoring exporter
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:         *projectID,
		MetricPrefix:      "opencensus-demo",
		ReportingInterval: 60 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create the Cloud Monitoring exporter: %v", err)
	}
	defer sd.Flush()
	// Start the metrics exporter
	sd.StartMetricsExporter()
	defer sd.StopMetricsExporter()
	// [END monitoring_sli_metrics_opencensus_exporter]
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	// create context
	ctx, _ := tag.New(context.Background())
	requestReceived := time.Now()
	// [START monitoring_sli_metrics_opencensus_counts]
	// count the request
	stats.Record(ctx, requestCount.M(1))

	// fail 10% of the time
	if rand.Intn(100) > 90 {
		// count the failed request
		stats.Record(ctx, failedRequestCount.M(1))
		// [END monitoring_sli_metrics_opencensus_counts]
		fmt.Fprintf(w, "intentional error!")
		// record latency for failure
		// [START monitoring_sli_metrics_opencensus_latency]
		stats.Record(ctx, responseLatency.M(time.Since(requestReceived).Seconds()))
		// [END monitoring_sli_metrics_opencensus_latency]
		return
	} else {
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(delay)
		fmt.Fprintf(w, "Succeeded after "+strconv.Itoa(int(delay)/1000000)+" ms")
		// record latency for success
		stats.Record(ctx, responseLatency.M(time.Since(requestReceived).Seconds()))
		return
	}
}
