package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("please tell how much you want")
		os.Exit(1)
	}
	val, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
