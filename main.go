package main

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
	mem_tracking   []chan bool
)

//Mem functions
func (m *memEater) memGetHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(mem_tracking) == 1 {
		io.WriteString(res, "memEater running!")
	} else {
		io.WriteString(res, "memEater not running!")
	}
}

func (m *memEater) cleanUpMemory(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing mem")
	for _, child := range mem_tracking {
		close(child)
	}

	mem_tracking = nil
	io.WriteString(res, "mem released!")
	C.free(unsafe.Pointer(m.echoOut))
	debug.FreeOSMemory()

	res.WriteHeader(http.StatusAccepted)
}

func (m *memEater) memPutHandler(res http.ResponseWriter, r *http.Request) {
	val, _ := mux.Vars(r)["val"]

	iVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}

	if len(mem_tracking) == 1 {
		io.WriteString(res, "memEater running!")
		return
	}

	signal := make(chan bool)

	go memEaterJob(m, iVal, signal)

	mem_tracking = append(mem_tracking, signal)

	io.WriteString(res, "memEater started!")
	res.WriteHeader(http.StatusAccepted)
}

type memEater struct{
	echoOut *C.char
}

func  memEaterJob(m *memEater, val int, signal chan bool) {

	cVal := C.int(val)
	m.echoOut = C.cEater(cVal)

	for s := range signal {
		fmt.Println(s)
	}
}

func main() {
	r := mux.NewRouter()
	m := memEater{}
	c := cpuBurner{}

	//mem handlers
	r.HandleFunc("/memeater", m.memGetHandler).Methods("GET")
	r.HandleFunc("/memeater/start/{val}", m.memPutHandler).Methods("GET")
	r.HandleFunc("/memeater/free", m.cleanUpMemory).Methods("GET")

	//CPU handlers
	r.HandleFunc("/cpuburner", c.cpuBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/start", c.cpuStartHandler).Methods("GET")
	r.HandleFunc("/cpuburner/stop", c.cpuStopHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
