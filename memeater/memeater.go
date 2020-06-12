package memeater

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//char* cEater(int s);
import "C"

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"unsafe"
)

// MemGetHandler returns memEater state
func (m *MemEaterHandler) MemGetHandler(res http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["id"]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, "{'error': 'id not found'}")
		return
	}
	if ms, ok := m.MemEater[id]; ok {
		data, err := json.Marshal(ms)
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

//CleanUpMemory stops memEaterJob and frees memory
func (m *MemEaterHandler) CleanUpMemory(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing mem")
	id, found := mux.Vars(r)["id"]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		res.Header().Set("Content-Type", "application/json")
		io.WriteString(res, "{'error': 'id not found'}")
		return
	}
	if ms, ok := m.MemEater[id]; ok {
		C.free(unsafe.Pointer(ms.echoOut))
		debug.FreeOSMemory()
		delete(m.MemEater, id)
		res.WriteHeader(http.StatusNoContent)
		return
	}
	res.WriteHeader(http.StatusNotFound)
	res.Header().Set("Content-Type", "application/json")
	io.WriteString(res, "{'error': 'id not found'}")
}

// MemPutHandler starts memEaterJob receives a quantity of memory in mb and time
func (m *MemEaterHandler) MemPutHandler(res http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mem memParams
	err := decoder.Decode(&mem)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	ms := &MemEater{
		Mem: mem.MemMb,
		ID:  xid.New().String(),
	}
	go memEaterJob(ms)
	m.MemEater[ms.ID] = ms
	data, err := json.Marshal(*ms)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// MemEater struct where attibutes are managed
type MemEater struct {
	echoOut *C.char
	Mem     int
	ID      string
}

// MemEaterHandler for handling different MemEater
type MemEaterHandler struct {
	MemEater map[string]*MemEater
}

type memParams struct {
	MemMb int
}

// NewMemEaterHandler creates a new map of MemEaters
func NewMemEaterHandler() *MemEaterHandler {

	return &MemEaterHandler{
		MemEater: map[string]*MemEater{},
	}

}

func memEaterJob(m *MemEater) {
	m.echoOut = C.cEater(C.int(m.Mem))
}
