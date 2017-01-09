package main

import (
	"flag"
	"github.com/Symantec/tricorder/go/tricorder"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"log"
	"math"
	"net/http"
	"net/rpc"
	"time"
)

var (
	kStartTime = time.Now()
)

func elapsed() float64 {
	return float64(time.Since(kStartTime)) / float64(time.Second)
}

func main() {
	flag.Parse()
	rpc.HandleHTTP()
	tricorder.UnregisterPath("/proc")
	if err := tricorder.RegisterMetric(
		"testing/increasing",
		elapsed,
		units.None,
		"Ever increasing"); err != nil {
		log.Fatal(err)
	}
	if err := tricorder.RegisterMetric(
		"testing/decreasing",
		func() float64 {
			return 1000000.0 - math.Mod(elapsed(), 1000000.0)
		},
		units.None,
		"Decreasing"); err != nil {
		log.Fatal(err)
	}
	if err := tricorder.RegisterMetric(
		"testing/sawtooth",
		func() float64 {
			x := elapsed()
			if x < 600.0/2.0 {
				return x
			}
			return 600.0 - x
		},
		units.None,
		"Saw tooth"); err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":6910", nil); err != nil {
		log.Fatal(err)
	}
}
