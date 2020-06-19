package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/JavierClairvaux/px-usewithcare/util"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Burner interface {
	Start()
	Stop()
	IsRunning() bool
	ID() uuid.UUID
}

type NewBurner func(io.ReadCloser) (Burner, error)

type BurnerHandler struct {
	instances map[uuid.UUID]Burner
	NewBurner NewBurner
	mutex     sync.Mutex
}

func NewBurnerHandler(f NewBurner) *BurnerHandler {

	c := &BurnerHandler{
		instances: make(map[uuid.UUID]Burner),
		mutex:     sync.Mutex{},
		NewBurner: f,
	}
	c.MonitorStoppedJobs()
	return c
}

// MonitorStoppedJobs monitors and cleans stopped jobs
func (c *BurnerHandler) MonitorStoppedJobs() {
	go removeJobs(c)
}

func removeJobs(c *BurnerHandler) {
	for {
		time.Sleep(1 * time.Second)
		for _, cs := range c.instances {
			if !cs.IsRunning() {
				c.mutex.Lock()
				log.Printf("Deleting Job with ID %s", cs.ID().String())
				delete(c.instances, cs.ID())
				c.mutex.Unlock()
			}
		}
	}
}

// StartHandler HTTP handler that starts Burner
func (c *BurnerHandler) StartHandler(w http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	burner, err := c.NewBurner(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go burner.Start()
	c.instances[burner.ID()] = burner
	data, err := json.Marshal(burner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// StopHandler HTTP handler that stops Burner
func (c *BurnerHandler) StopHandler(w http.ResponseWriter, r *http.Request) {
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
	if cs, ok := c.instances[id]; ok {
		cs.Stop()
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

// GetHandler HTTP handler that returns the Burner state
func (c *BurnerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("Invalid UUID: %s", err)
		return
	}
	if cs, ok := c.instances[id]; ok {
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

// ListHandler returns a list of active CPU burners
func (c *BurnerHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	log.Println("Listing active burners")
	s := []Burner{}
	for _, cs := range c.instances {
		s = append(s, cs)
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Cannot serialize list %s", err.Error())
	}
	w.Write(data)
}

func HandlePaths(path string, h *BurnerHandler, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), h.GetHandler).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s", path), h.StartHandler).Methods("POST")
	// this function would need to a pluralizer, but I don't think is worth to add it for now
	r.HandleFunc(fmt.Sprintf("/%ss", path), h.ListHandler).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), h.StopHandler).Methods("DELETE")
}
