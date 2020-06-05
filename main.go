package main

import (
        "github.com/gorilla/mux"
        "io"
         "net/http"
	 "bytes"
	 "log"
	 "time"
	 "strconv"
	 "runtime"
 )

func memGetHandler(res http.ResponseWriter, r *http.Request){
         //vars := mux.Vars(r)
         res.WriteHeader(http.StatusOK)
         //io.WriteString(res, "dog dog dog")
         io.WriteString(res, strconv.Itoa(runtime.NumGoroutine()))
}

func memPutHandler(res http.ResponseWriter, r *http.Request){
        //vars := mux.Vars(r)
        val, _:= mux.Vars(r)["val"]
        //res.WriteHeader(http.StatusOK)
        io.WriteString(res, val)
	memint, _ := strconv.Atoi(val)
	go memEater( memint )
        res.WriteHeader(http.StatusAccepted)
}

func memEater(val int) {

	const n = 1000000
	b1 := make([]byte, n, n)
	b2 := make([]byte, n, n)
	b3 := make([]byte, n, n)

	w1, w2, w3 := bytes.NewBuffer(b1), bytes.NewBuffer(b2), bytes.NewBuffer(b3)
	w := io.MultiWriter(w1, w2, w3)

	for i := 0; i < val; i++ {
		w.Write(make([]byte, n, n))
	}

	for {
		time.Sleep(30 * time.Second)
		log.Println("eater happy ..")
	}
}

func main() {
        r := mux.NewRouter()
        r.HandleFunc("/memeater", memGetHandler).Methods("GET")
        r.HandleFunc("/memeater/{val}", memPutHandler).Methods("PUT")
        http.Handle("/", r)
        http.ListenAndServe(":8080", r)
}
