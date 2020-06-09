package cpuburner

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
	cpuTracking   []chan bool
)

func cpuBurn(c *CPUBurner, signal chan bool) {
	for {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
		if !c.running{
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

func cpuBurnerJob(c *CPUBurner, signal chan bool) {
	runtime.GOMAXPROCS(numBurn)
	fmt.Printf("Burning %d CPUs/cores\n", numBurn)
	for i := 0; i < numBurn; i++ {
		go cpuBurn(c, signal)
	}
	if updateInterval > 0 {
		t := time.Tick(time.Duration(updateInterval) * time.Second)
		for secs := updateInterval; ; secs += updateInterval {
			if !c.running{
				return
			}
			<-t
			fmt.Printf("%d seconds\n", secs)
		}
	} else {
		select {} // wait forever
	}
}

type CPUBurner struct{
	running bool
}

//CPU burners handlers
func (c *CPUBurner) CPUBurnerHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(cpuTracking) == 1 {
		io.WriteString(res, "{'cpuBurner': 'started'}")
	} else {
		io.WriteString(res, "{'cpuBurner': 'stopped'}")
	}
}

func (c *CPUBurner) CPUStartHandler(res http.ResponseWriter, r *http.Request) {
	if len(cpuTracking) == 1 {
		io.WriteString(res, "{'cpuBurner': 'started'}")
		return
	}

	signal := make(chan bool)

	c.running = true
	go cpuBurnerJob(c, signal)
	cpuTracking = append(cpuTracking, signal)


	io.WriteString(res, "{'cpuBurner': 'started'}")
}

func (c *CPUBurner) CPUStopHandler(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing CPU")
	c.running = false
	for _, child := range cpuTracking {
		close(child)
	}

	cpuTracking = nil
	io.WriteString(res, "{'cpuBurner': 'stopped'}")
	res.WriteHeader(http.StatusAccepted)
}
