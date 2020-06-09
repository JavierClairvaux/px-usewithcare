package memeater

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//char* cEater(int s);
import "C"

import (
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"fmt"
	"unsafe"
	"github.com/gorilla/mux"
)

var (
	memTracking   []chan bool
)

//Mem functions
func (m *MemEater) MemGetHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(memTracking) == 1 {
		io.WriteString(res, "{'memEater': 'started'}")
	} else {
		io.WriteString(res, "{'memEater': 'stopped'}")
	}
}

func (m *MemEater) CleanUpMemory(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing mem")
	for _, child := range memTracking {
		close(child)
	}
	memTracking = nil
	io.WriteString(res, "{'memEater': 'stopped'}")
	C.free(unsafe.Pointer(m.echoOut))
	debug.FreeOSMemory()
	res.WriteHeader(http.StatusAccepted)
}

func (m *MemEater) MemPutHandler(res http.ResponseWriter, r *http.Request) {
	val, _ := mux.Vars(r)["val"]
	iVal, err := strconv.Atoi(val)

	if err != nil {
		log.Fatal(err)
	}
	if len(memTracking) == 1 {
		io.WriteString(res, "{'memEater': 'started'}")
		return
	}
	signal := make(chan bool)
	go memEaterJob(m, iVal, signal)
	memTracking = append(memTracking, signal)
	io.WriteString(res, "{'memEater': 'started'}")
	res.WriteHeader(http.StatusAccepted)
}

type MemEater struct{
	echoOut *C.char
}

func  memEaterJob(m *MemEater, val int, signal chan bool) {
	cVal := C.int(val)
	m.echoOut = C.cEater(cVal)

	for s := range signal {
		fmt.Println(s)
	}
}
