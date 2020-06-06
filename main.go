package main

import (
	"io"
	"log"
	"net/http"
	"os"
	//"runtime"
	"strconv"
	"syscall"
	//"flag"
	//"fmt"
	//"time"

	"github.com/gorilla/mux"
)

var (
	mem_tracking 	[]int
	cpu_tracking 	[]int
	numBurn		int
	updateInterval 	int
)
//Mem functions
func memGetHandler(res http.ResponseWriter, r *http.Request) {
	res.WriteHeader(http.StatusOK)
	if len(mem_tracking) == 1 {
		io.WriteString(res, "memEater running!")
	} else {
		io.WriteString(res, "memEater not running!")
	}
}

func cleanUpMemory(res http.ResponseWriter, r *http.Request) {
	log.Println("Releasing mem")
	for _, child := range mem_tracking {
		syscall.Kill(child, syscall.SIGQUIT)
	}

	mem_tracking = nil
	
	io.WriteString(res, "mem released!")

	res.WriteHeader(http.StatusAccepted)
}

func memPutHandler(res http.ResponseWriter, r *http.Request) {
	val, _ := mux.Vars(r)["val"]

	_, err := strconv.Atoi(val)

	if err != nil {
		log.Fatal(err)
	}


	if len(mem_tracking) == 1 {
		io.WriteString(res, "memEater running!")
		return
	}

	binary := "./bin/memeater"
	childPID, _ := syscall.ForkExec(binary, []string{binary, val}, &syscall.ProcAttr{
		Dir: "./",
		Env: os.Environ(),
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
		Files: []uintptr{0, 1, 2}, // print message to the same pty
	})
	log.Printf("child %d", childPID)

	if childPID != 0 {
		mem_tracking = append(mem_tracking, childPID)
	}

	io.WriteString(res, "memEater started!")
}

//CPU burners functions
func cpuBurnerHandler(res http.ResponseWriter, r *http.Request){
	res.WriteHeader(http.StatusOK)
	if len(cpu_tracking) == 1 {
		io.WriteString(res, "cpuBurner running!")
	} else {
		io.WriteString(res, "cpuBurner not running!")
	}

}

func cpuStartHandler(res http.ResponseWriter, r *http.Request){
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


func cpuStopHandler(res http.ResponseWriter, r *http.Request){
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

	//mem handlers 
	r.HandleFunc("/memeater", memGetHandler).Methods("GET")
	r.HandleFunc("/memeater/{val}", memPutHandler).Methods("PUT")
	r.HandleFunc("/memeater/free", cleanUpMemory).Methods("GET")

	//CPU handlers
	r.HandleFunc("/cpuburner", cpuBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/start", cpuStartHandler).Methods("PUT")
	r.HandleFunc("/cpuburner/stop", cpuStopHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
