package main

import (
	"net/http"

	"github.com/JavierClairvaux/px-usewithcare/cpu"
	memory "github.com/JavierClairvaux/px-usewithcare/memory"

	"github.com/JavierClairvaux/px-usewithcare/handler"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	c := handler.NewBurnerHandler(cpu.NewBurner)
	m := handler.NewBurnerHandler(memory.NewBurner)

	// handle cpu burner
	handler.HandlePaths("cpuburner", c, r)

	//handle memeater
	handler.HandlePaths("memeater", m, r)

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
