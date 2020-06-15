package cpuburner

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func cpuBurn(c *CPUBurner) {
	for {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
		if !c.Running {
			break
		}
	}
}

func cpuBurnerJob(c *CPUBurner) {
	fmt.Printf("Burning %d CPUs/cores\n", c.NumBurn)
	for i := 0; i < c.NumBurn; i++ {
		go cpuBurn(c)
	}
	fmt.Printf("Sleeping %d miliseconds\n", c.TTL)
	for start := time.Now(); time.Since(start) < time.Millisecond*time.Duration(c.TTL); {
	}
	c.NumBurn = 0
	c.Running = false
	c.TTL = 0
}

// CPUBurner struct where all the parameters are stored
type CPUBurner struct {
	Running bool   `json:"running"`
	NumBurn int    `json:"num_burn"`
	TTL     int    `json:"ttl"`
	ID      string `json:"id,omitempty"`
}

type cParams struct {
	Count int
	TTL   int
}

// CPUBurnerHandler map for managing processes
type cpuBurnerHandler struct {
	mutex     sync.Mutex
	cpuBurner map[string]*CPUBurner
}

// NewCPUBurnerHandler Handler that returns cpuBurnerHandler with new ID
func NewCPUBurnerHandler() *cpuBurnerHandler {

	return &cpuBurnerHandler{
		cpuBurner: make(map[string]*CPUBurner),
		mutex:     sync.Mutex{},
	}
}

// RemoveStoppedJobs stopped jobs cleanner
func (c *cpuBurnerHandler) RemoveStoppedJobs() {
	go removeJobs(c)
}

func removeJobs(c *cpuBurnerHandler) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	for {
		time.Sleep(1 * time.Second)
		for _, cs := range c.cpuBurner {
			if !cs.Running {
				delete(c.cpuBurner, cs.ID)
			}
		}
	}
}

// CPUBurnerHandler HTTP handler that returns cpuBurner state
func (c *cpuBurnerHandler) CPUBurnerHandler(res http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	//key, ok := r.URL.Query()["id"]
	id, found := mux.Vars(r)["id"]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, "{'error': 'id not found'}")
		return
	}
	if cs, ok := c.cpuBurner[id]; ok {
		data, err := json.Marshal(cs)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		return
	}
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, "{'error': 'id not found'}")
}

// CPUStartHandler HTTP handler that starts cpuBurnerJob
func (c *cpuBurnerHandler) CPUStartHandler(res http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	decoder := json.NewDecoder(r.Body)
	var p cParams
	err := decoder.Decode(&p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	cs := &CPUBurner{
		Running: true,
		NumBurn: p.Count,
		TTL:     p.TTL,
		ID:      xid.New().String(),
	}
	go cpuBurnerJob(cs)
	c.cpuBurner[cs.ID] = cs
	data, err := json.Marshal(*cs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// CPUStopHandler HTTP handler that stops cpuBurnerJob
func (c *cpuBurnerHandler) CPUStopHandler(res http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	log.Println("Releasing CPU")
	id, found := mux.Vars(r)["id"]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, "{'error': 'id not found'}")
		return
	}
	if cs, ok := c.cpuBurner[id]; ok {
		cs.Running = false
		cs.NumBurn = 0
		cs.TTL = 0
		delete(c.cpuBurner, id)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, "{'error': 'id not found'}")
}
