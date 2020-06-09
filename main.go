package main

import (
	"net/http"
	"github.com/JavierClairvaux/px-usewithcare/cpuburner"
	"github.com/JavierClairvaux/px-usewithcare/memeater"

	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	m := memeater.MemEater{}
	c := cpuburner.CPUBurner{}

	//mem handlers
	r.HandleFunc("/memeater", m.MemGetHandler).Methods("GET")
	r.HandleFunc("/memeater/start/{val}", m.MemPutHandler).Methods("GET")
	r.HandleFunc("/memeater/free", m.CleanUpMemory).Methods("GET")

	//CPU handlers
	r.HandleFunc("/cpuburner", c.CPUBurnerHandler).Methods("GET")
	r.HandleFunc("/cpuburner/start", c.CPUStartHandler).Methods("GET")
	r.HandleFunc("/cpuburner/stop", c.CPUStopHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
