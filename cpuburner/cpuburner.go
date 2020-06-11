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
	if c.UpdateInterval > 0 {
		fmt.Printf("Sleeping %d miliseconds\n", c.UpdateInterval)
		for start := time.Now(); time.Since(start) < time.Millisecond*time.Duration(c.UpdateInterval); {
		}
		c.NumBurn = 0
		c.Running = false
		c.UpdateInterval = 0
	} else {
		select {} // wait forever
	}
}

type CPUBurner struct {
	Running        bool   `json:"running"`
	NumBurn        int    `json:"num_burn"`
	UpdateInterval int    `json:"ttl"`
	ID             string `json:"id,omitempty"`
}

type cParams struct {
	Count int
	TTL   int
}

//CPUBurnerHandler map for managing processes
type CPUBurnerHandler struct {
	cpuBurner map[string]*CPUBurner
	id        xid.ID
}

//NewCPUBurnerHandler Handler that returns CPUBurnerHandler with new ID
func NewCPUBurnerHandler() *CPUBurnerHandler {
	guid := xid.New()

	return &CPUBurnerHandler{
		cpuBurner: map[string]*CPUBurner{},
		id:        guid,
	}
}

//RemoveStoppedJobs stopped jobs cleanner
func (c *CPUBurnerHandler) RemoveStoppedJobs() {
	go removeJobs(c)
}

func removeJobs(c *CPUBurnerHandler) {
	for {
		time.Sleep(1 * time.Second)
		for _, cs := range c.cpuBurner {
			if !cs.Running {
				delete(c.cpuBurner, cs.ID)
			}
		}
	}
}

//CPU burners handlers
//CPUBurnerHandler HTTP handler that returns cpuBurner state
func (c *CPUBurnerHandler) CPUBurnerHandler(res http.ResponseWriter, r *http.Request) {
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

//CPUStartHandler HTTP handler that starts cpuBurnerJob
func (c *CPUBurnerHandler) CPUStartHandler(res http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p cParams
	err := decoder.Decode(&p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	cs := &CPUBurner{
		Running:        true,
		NumBurn:        p.Count,
		UpdateInterval: p.TTL,
		ID:             c.id.String(),
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

//CPUStopHandler HTTP handler that stops cpuBurnerJob
func (c *CPUBurnerHandler) CPUStopHandler(res http.ResponseWriter, r *http.Request) {
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
		cs.UpdateInterval = 0
		delete(c.cpuBurner, id)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, "{'error': 'id not found'}")
}
