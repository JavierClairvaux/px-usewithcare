package main

import (
	"github.com/JavierClairvaux/px-usewithcare/cpuburner"
	"github.com/JavierClairvaux/px-usewithcare/memeater"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	m := memeater.MemEater{}
	c := cpuburner.NewCPUBurnerHandler()
	c.RemoveStoppedJobs()

	//mem handlers
	r.HandleFunc("/memeater", m.MemGetHandler).Methods("GET")
	r.HandleFunc("/memeater/start/{val}", m.MemPutHandler).Methods("GET")
	r.HandleFunc("/memeater/free", m.CleanUpMemory).Methods("GET")

	//CPU handlers
	r.HandleFunc("/cpuburner/{id}", c.CPUBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/start", c.CPUStartHandler).Methods("POST")
	r.HandleFunc("/cpuburner/stop/{id}", c.CPUStopHandler).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
