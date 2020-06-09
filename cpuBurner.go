package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"net/http"
	"io"
	"log"
)

var (
	numBurn        int
	updateInterval int
	cpu_tracking   []chan bool
)

func cpuBurn(c *cpuBurner, signal chan bool) {
	for {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
		if !c.Running{
			break
		}
	}
}

func init() {
	flag.IntVar(&numBurn, "n", 0, "number of cores to burn (0 = all)")
	flag.IntVar(&updateInterval, "u", 10, "seconds between updates (0 = don't update)")
	flag.Parse()
	if numBurn <= 0 {
		numBurn = runtime.NumCPU()
	}
}

func cpuBurnerJob(c *cpuBurner, signal chan bool) {
	runtime.GOMAXPROCS(numBurn)
	fmt.Printf("Burning %d CPUs/cores\n", numBurn)
	for i := 0; i < numBurn; i++ {
		go cpuBurn(c, signal)
	}
	if updateInterval > 0 {
		t := time.Tick(time.Duration(updateInterval) * time.Second)
		for secs := updateInterval; ; secs += updateInterval {
			if !c.Running{
				return
			}
			<-t
			fmt.Printf("%d seconds\n", secs)
		}
	} else {
		select {} // wait forever
	}
}

type cpuBurner struct{
	Running bool
}

//CPU burners handlers
func (c *cpuBurner) cpuBurnerHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(cpu_tracking) == 1 {
		io.WriteString(res, "cpuBurner running!")
	} else {
		io.WriteString(res, "cpuBurner not running!")
	}
}

func (c *cpuBurner) cpuStartHandler(res http.ResponseWriter, r *http.Request) {
	if len(cpu_tracking) == 1 {
		io.WriteString(res, "cpuBurner running!")
		return
	}

	signal := make(chan bool)

	c.Running = true
	go cpuBurnerJob(c, signal)
	cpu_tracking = append(cpu_tracking, signal)


	io.WriteString(res, "cpuBurner started!")
}

func (c *cpuBurner) cpuStopHandler(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing CPU")
	c.Running = false
	for _, child := range cpu_tracking {
		close(child)
	}

	cpu_tracking = nil
	io.WriteString(res, "cpu released")
	res.WriteHeader(http.StatusAccepted)
}
