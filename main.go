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
	"os"
	"runtime/debug"
	"strconv"
	"syscall"
	"fmt"
	"unsafe"

	"github.com/gorilla/mux"
)

var (
	mem_tracking   []chan bool
	cpu_tracking   []int
	numBurn        int
	updateInterval int
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

//CPU burners functions
func cpuBurnerHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(cpu_tracking) == 1 {
		io.WriteString(res, "cpuBurner running!")
	} else {
		io.WriteString(res, "cpuBurner not running!")
	}

}

func cpuStartHandler(res http.ResponseWriter, r *http.Request) {
	if len(cpu_tracking) == 1 {
		io.WriteString(res, "cpuBurner running!")
		return
	}

	binary := "./bin/cpuburner"
	childPID, _ := syscall.ForkExec(binary, []string{binary}, &syscall.ProcAttr{
		Dir: "./",
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
		Files: []uintptr{0, 1, 2}, // print message to the same pty
	})
	log.Printf("child %d", childPID)

	if childPID != 0 {
		cpu_tracking = append(cpu_tracking, childPID)
	}

	io.WriteString(res, "cpuBurner started!")
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

func cpuStopHandler(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing CPU")
	for _, child := range cpu_tracking {
		syscall.Kill(child, syscall.SIGQUIT)
	}

	cpu_tracking = nil

	io.WriteString(res, "cpu released")

	res.WriteHeader(http.StatusAccepted)
}
func main() {
	r := mux.NewRouter()
	m:= memEater{}

	//mem handlers
	r.HandleFunc("/memeater", m.memGetHandler).Methods("GET")
	r.HandleFunc("/memeater/start/{val}", m.memPutHandler).Methods("GET")
	r.HandleFunc("/memeater/free", m.cleanUpMemory).Methods("GET")

	//CPU handlers
	r.HandleFunc("/cpuburner", cpuBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/start", cpuStartHandler).Methods("PUT")
	r.HandleFunc("/cpuburner/stop", cpuStopHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
