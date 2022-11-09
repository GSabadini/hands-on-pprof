package main

import (
	"context"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
	"time"
)

func main() {
	//f, err := os.Create("trace.out")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//
	//err = trace.Start(f)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer trace.Stop()

	err := statsviz.RegisterDefault()
	if err != nil {
		log.Fatal(err)
	}

	initMeter()
	//initTracer()
	shutdown := NewTracer()
	defer shutdown()

	db := NewMySQLConnection()
	defer db.Close()
	tracer := otel.GetTracerProvider()

	labels := pprof.Labels("goroutine", "test")
	pprof.Do(context.Background(), labels, func(ctx context.Context) {
		for x := 1; x < 10; x++ {
			go GoroutineWithLabel()
		}
	})

	http.HandleFunc("/select", SelectHandler(1, db, tracer))
	http.HandleFunc("/insert", InsertHandler(randomString(), db, tracer))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":6061", nil))
}

func GoroutineWithLabel() {
	fmt.Println("ehehe")
	time.Sleep(1 * time.Minute)
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyz"
	c := charset[rand.Intn(len(charset))]

	return string(c)
}
