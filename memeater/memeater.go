package memeater

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//char* cEater(int s);
import "C"

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"unsafe"
)

//Mem functions
//MemGetHandler returns memEater state
func (m *MemEater) MemGetHandler(res http.ResponseWriter, r *http.Request) {
	if len(m.memTracking) == 1 {
		io.WriteString(res, "{'memEater': 'started'}")
	} else {
		io.WriteString(res, "{'memEater': 'stopped'}")
	}
}

//CleanUpMemory stops memEaterJob and frees memory
func (m *MemEater) CleanUpMemory(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing mem")
	for _, child := range m.memTracking {
		close(child)
	}
	m.memTracking = nil
	io.WriteString(res, "{'memEater': 'stopped'}")
	C.free(unsafe.Pointer(m.echoOut))
	debug.FreeOSMemory()
	res.WriteHeader(http.StatusAccepted)
}

//MemPutHandler starts memEaterJob receives a quantity of memory in mb and time
func (m *MemEater) MemPutHandler(res http.ResponseWriter, r *http.Request) {
	val, _ := mux.Vars(r)["val"]
	iVal, err := strconv.Atoi(val)

	if err != nil {
		log.Fatal(err)
	}
	if len(m.memTracking) == 1 {
		io.WriteString(res, "{'memEater': 'started'}")
		return
	}
	signal := make(chan bool)
	go memEaterJob(m, iVal, signal)
	m.memTracking = append(m.memTracking, signal)
	io.WriteString(res, "{'memEater': 'started'}")
	res.WriteHeader(http.StatusAccepted)
}

type MemEater struct {
	echoOut     *C.char
	memTracking []chan bool
}

func memEaterJob(m *MemEater, val int, signal chan bool) {
	cVal := C.int(val)
	m.echoOut = C.cEater(cVal)
}
