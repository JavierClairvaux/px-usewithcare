package memeater

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//char* cEater(int s);
import "C"

import (
	"encoding/json"
	"fmt"
	"github.com/JavierClairvaux/px-usewithcare/util"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"unsafe"
)

// MemGetHandler returns memEater state
func (m *memEaterHandler) MemGetHandler(w http.ResponseWriter, r *http.Request) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
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
		fmt.Printf("Something wend wrong: %s", err)
		return
	}
	if ms, ok := m.MemEater[id]; ok {
		data, err := json.Marshal(ms)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

//CleanUpMemory stops memEaterJob and frees memory
func (m *memEaterHandler) CleanUpMemory(w http.ResponseWriter, r *http.Request) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	log.Println("Releasing mem")
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
	if ms, ok := m.MemEater[id]; ok {
		C.free(unsafe.Pointer(ms.echoOut))
		debug.FreeOSMemory()
		delete(m.MemEater, id)
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

// MemPutHandler starts memEaterJob receives a quantity of memory in mb and time
func (m *memEaterHandler) MemPutHandler(w http.ResponseWriter, r *http.Request) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	decoder := json.NewDecoder(r.Body)
	var mem memParams
	err := decoder.Decode(&mem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ms := &MemEater{
		Mem: mem.MemMb,
		ID:  uuid.NewV4(),
	}
	go memEaterJob(ms)
	m.MemEater[ms.ID] = ms
	data, err := json.Marshal(*ms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// MemEater struct where attibutes are managed
type MemEater struct {
	echoOut *C.char
	Mem     int       `json:"mem_mb"`
	ID      uuid.UUID `json:"id",omitempty`
}

// MemEaterHandler for handling different MemEater
type memEaterHandler struct {
	MemEater map[uuid.UUID]*MemEater
	mutex    sync.Mutex
}

type memParams struct {
	MemMb int
}

// NewMemEaterHandler creates a new map of MemEaters
func NewMemEaterHandler() *memEaterHandler {

	return &memEaterHandler{
		MemEater: map[uuid.UUID]*MemEater{},
		mutex:    sync.Mutex{},
	}

}

func memEaterJob(m *MemEater) {
	m.echoOut = C.cEater(C.int(m.Mem))
}
