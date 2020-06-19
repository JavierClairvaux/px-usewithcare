package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/JavierClairvaux/px-usewithcare/util"
	"github.com/gertd/go-pluralize"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Burner is an interface that defines the behaviour of the burners
type Burner interface {
	Start()
	Stop()
	ID() uuid.UUID
}

// NewBurner is a type function to return a new instance of a Burner
type NewBurner func(io.ReadCloser) (Burner, error)

// BurnerHandler is an HTTP Handler that manages a Burner
type BurnerHandler struct {
	instances *sync.Map
	NewBurner NewBurner
	delete    chan uuid.UUID
}

// NewBurnerHandler returns a new instance of a BurnerHandler
func NewBurnerHandler(f NewBurner) *BurnerHandler {

	c := &BurnerHandler{
		instances: &sync.Map{},
		NewBurner: f,
	}
	c.delete = make(chan uuid.UUID)
	return c
}

// StartHandler HTTP handler that starts Burner
func (c *BurnerHandler) StartHandler(w http.ResponseWriter, r *http.Request) {

	burner, err := c.NewBurner(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, loaded := c.instances.LoadOrStore(burner.ID(), burner)
	if loaded {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		data, err := util.GetHTTPError("could not store the burner")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	go burner.Start()

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

	idRaw, found := mux.Vars(r)["id"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		data, err := util.GetHTTPError("ID not found")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	id, err := uuid.FromString(idRaw)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	if cs, ok := c.instances.Load(id); ok {
		b := cs.(Burner)
		b.Stop()
		c.instances.Delete(b)
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
	if cs, ok := c.instances.Load(id); ok {
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

	log.Println("Listing active burners")
	s := []Burner{}

	c.instances.Range(func(key interface{}, value interface{}) bool {
		s = append(s, value.(Burner))
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Cannot serialize list %s", err.Error())
	}
	w.Write(data)
}

// HandlePaths sets the paths for each of the functions of a BurnerHandler
func HandlePaths(path string, h *BurnerHandler, r *mux.Router) {
	pluralize := pluralize.NewClient()

	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), h.GetHandler).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s", path), h.StartHandler).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s", pluralize.Plural(path)), h.ListHandler).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{id}", path), h.StopHandler).Methods("DELETE")
}
