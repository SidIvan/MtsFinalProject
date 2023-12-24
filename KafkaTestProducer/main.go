package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var async = flag.Bool("a", false, "use async")

//func main() {
//	flag.Parse()
//
//	ctx := context.Background()
//
//	logger := log.Default()
//
//	writer := kafka.NewWriter(kafka.WriterConfig{
//		Brokers:     []string{"127.0.0.1:29092"},
//		Topic:       "read_topic",
//		Async:       *async,
//		Logger:      kafka.LoggerFunc(logger.Printf),
//		ErrorLogger: kafka.LoggerFunc(logger.Printf),
//		BatchSize:   2000,
//	})
//	defer writer.Close()
//	testMsgpref := "{\n    \"id\": \"284655d6-0190-49e7-34e9-9b4060acc261\",\n    \"source\": \"/trip\",\n    \"type\": \"trip.event.created\",\n    \"datacontenttype\": \"application/json\",\n    \"time\": \"2023-11-09T17:31:00Z\",\n    \"data\": {\n        \"trip_id\": \""
//	testMsgpostf := "e82c42d6-b86f-4e2a-93a2-858413acb148\",\n        \"offer_id\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InN0cmluZyIsImZyb20iOnsibGF0IjowLCJsbmciOjB9LCJ0byI6eyJsYXQiOjAsImxuZyI6MH0sImNsaWVudF9pZCI6InN0cmluZyIsInByaWNlIjp7ImFtb3VudCI6OTkuOTUsImN1cnJlbmN5IjoiUlVCIn19.fg0Bv2ONjT4r8OgFqJ2tpv67ar7pUih2LhDRCRhWW3c\",\n        \"price\": {\n            \"currency\": \"RUB\",\n            \"amount\": 100\n        },\n        \"status\": \"DRIVER_SEARCH\",\n        \"from\": {\n            \"lat\": 0,\n            \"lng\": 0\n        },\n        \"to\": {\n            \"lat\": 0,\n            \"lng\": 0\n        }\n    }\n}"
//
//	for i := 0; i < 524_288; i++ {
//		var a int
//		fmt.Scanf("%d", &a)
//		err := writer.WriteMessages(ctx, kafka.Message{Key: []byte(strconv.Itoa(i)), Value: []byte(testMsgpref + strconv.Itoa(i) + testMsgpostf)})
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}

func main() {
	// Create a new counter metric using promauto
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "example_counter",
		Help: "This is an example counter",
	})

	// Start a goroutine to periodically increment the counter
	go func() {
		for {
			counter.Inc()
			time.Sleep(1 * time.Second)
		}
	}()

	// Expose the metrics endpoint via an HTTP server
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9090", nil)
}
