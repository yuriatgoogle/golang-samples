module github.com/GoogleCloudPlatform/golang-samples/monitoring

go 1.13

replace github.com/GoogleCloudPlatform/golang-samples => ./..

require (
	cloud.google.com/go v0.81.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.2
	github.com/prometheus/client_golang v1.10.0
	go.opencensus.io v0.23.0
	google.golang.org/api v0.44.0
	google.golang.org/genproto v0.0.0-20210415145412-64678f1ae2d5
)
