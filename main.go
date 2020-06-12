package main

import (
	"github.com/JavierClairvaux/px-usewithcare/cpuburner"
	"github.com/JavierClairvaux/px-usewithcare/memeater"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	m := memeater.NewMemEaterHandler()
	c := cpuburner.NewCPUBurnerHandler()
	c.RemoveStoppedJobs()

	//mem handlers
	r.HandleFunc("/memeater/{id}", m.MemGetHandler).Methods("GET")
	r.HandleFunc("/memeater", m.MemPutHandler).Methods("POST")
	r.HandleFunc("/memeater/{id}", m.CleanUpMemory).Methods("DELETE")

	//CPU handlers
	r.HandleFunc("/cpuburner/{id}", c.CPUBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/", c.CPUStartHandler).Methods("POST")
	r.HandleFunc("/cpuburner/{id}", c.CPUStopHandler).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
