package main

import (
	"encoding/json"
	"flag"
	"github.com/smithclay/synthetic-load-generator-go/generator"
	"github.com/smithclay/synthetic-load-generator-go/topology"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var paramsFile string
	var collectorUrl string
	var randSeed int64
	stdoutMode := false

	flag.StringVar(&paramsFile, "paramsFile", "REQUIRED", "topology JSON file")
	flag.StringVar(&collectorUrl, "collectorUrl", "", "URL to gRPC OpenTelemetry collector")
	flag.Int64Var(&randSeed, "randSeed", time.Now().UTC().UnixNano(), "random seed (int64)")

	flag.Parse()
	if collectorUrl == "" {
		stdoutMode = true
	}

	jsonFile, err := os.Open(paramsFile)
	if err != nil {
		log.Fatalf("could not open topology file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var file topology.File
	err = json.Unmarshal(byteValue, &file)
	if err != nil {
		log.Fatalf("could not parse topology file: %v", err)
	}
	traceGenerators := make([]*generator.ScheduledTraceGenerator, 0)
	for _, r := range file.RootRoutes {
		var tg *generator.ScheduledTraceGenerator

		tg = generator.NewScheduledTraceGenerator(file.Topology, r.Route, r.Service,
			generator.WithSeed(randSeed),
			generator.WithTracesPerHour(r.TracesPerHour),
			generator.WithGrpc(collectorUrl))

		if stdoutMode {
			tg = generator.NewScheduledTraceGenerator(file.Topology, r.Route, r.Service,
				generator.WithSeed(randSeed),
				generator.WithTracesPerHour(r.TracesPerHour))
		}

		traceGenerators = append(traceGenerators, tg)
	}

	for _, tg := range traceGenerators {
		go tg.Start()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down...")
	for _, tg := range traceGenerators {
		tg.Shutdown()
	}
	os.Exit(0)

}
