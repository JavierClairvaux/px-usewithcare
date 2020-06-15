package cpuburner

import (
	"encoding/json"
	"fmt"
	"github.com/JavierClairvaux/px-usewithcare/util"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
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
	fmt.Printf("Burning %d CPUs/cow\n", c.NumBurn)
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
	Running bool      `json:"running"`
	NumBurn int       `json:"num_burn"`
	TTL     int       `json:"ttl"`
	ID      uuid.UUID `json:"id,omitempty"`
}

type cParams struct {
	Count int
	TTL   int
}

// CPUBurnerHandler map for managing processes
type cpuBurnerHandler struct {
	mutex     sync.Mutex
	cpuBurner map[uuid.UUID]*CPUBurner
}

// NewCPUBurnerHandler Handler that returns cpuBurnerHandler with new ID
func NewCPUBurnerHandler() *cpuBurnerHandler {

	return &cpuBurnerHandler{
		cpuBurner: make(map[uuid.UUID]*CPUBurner),
		mutex:     sync.Mutex{},
	}
}

// RemoveStoppedJobs stopped jobs cleanner
func (c *cpuBurnerHandler) RemoveStoppedJobs() {
	go removeJobs(c)
}

func removeJobs(c *cpuBurnerHandler) {
	for {
		time.Sleep(1 * time.Second)
		for _, cs := range c.cpuBurner {
			if !cs.Running {
				c.mutex.Lock()
				delete(c.cpuBurner, cs.ID)
				c.mutex.Unlock()
			}
		}
	}
}

// CPUBurnerHandler HTTP handler that returns cpuBurner state
func (c *cpuBurnerHandler) CPUBurnerHandler(w http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	idRaw, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		data, err := util.GetHTTPError("ID not found")
		if err != nil {
			log.Fatalf("Cannot serialize error %s", err.Error())
		}
		w.Write(data)
		return
	}
	id, err := uuid.FromString(idRaw)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	if cs, ok := c.cpuBurner[id]; ok {
		data, err := json.Marshal(cs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			io.WriteString(w, "Failed to serialize output!")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	data, err := util.GetHTTPError("ID not found")
	if err != nil {
		log.Fatalf("Cannot serialize error %s", err.Error())
	}
	w.Write(data)
}

// CPUStartHandler HTTP handler that starts cpuBurnerJob
func (c *cpuBurnerHandler) CPUStartHandler(w http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	decoder := json.NewDecoder(r.Body)
	var p cParams
	err := decoder.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cs := &CPUBurner{
		Running: true,
		NumBurn: p.Count,
		TTL:     p.TTL,
		ID:      uuid.NewV4(),
	}
	go cpuBurnerJob(cs)
	c.cpuBurner[cs.ID] = cs
	data, err := json.Marshal(*cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// CPUStopHandler HTTP handler that stops cpuBurnerJob
func (c *cpuBurnerHandler) CPUStopHandler(w http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	log.Println("Releasing CPU")
	idRaw, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		data, err := util.GetHTTPError("ID not found")
		if err != nil {
			log.Fatalf("Cannot serialize error %s", err.Error())
		}
		w.Write(data)
		return
	}
	id, err := uuid.FromString(idRaw)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	if cs, ok := c.cpuBurner[id]; ok {
		cs.Running = false
		cs.NumBurn = 0
		cs.TTL = 0
		delete(c.cpuBurner, id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	data, err := util.GetHTTPError("ID not found")
	if err != nil {
		log.Fatalf("Cannot serialize error %s", err.Error())
	}
	w.Write(data)
}
