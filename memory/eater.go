package eater

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//char* cEater(int s);
import "C"

import (
	"encoding/json"
	"runtime/debug"
	"unsafe"

	"github.com/JavierClairvaux/px-usewithcare/handler"
	uuid "github.com/satori/go.uuid"

	"io"
)

type Eater struct {
	echoOut *C.char
	Running bool      `json:"running"`
	Mem     int       `json:"mem_mb"`
	UUID    uuid.UUID `json:"id,omitempty"`
}

func NewBurner(body io.ReadCloser) (handler.Burner, error) {
	e := Eater{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&e)
	e.UUID = uuid.NewV4()
	e.Running = true
	return &e, err
}

func (e Eater) ID() uuid.UUID {
	return e.UUID
}

func (e Eater) IsRunning() bool {
	return e.Running
}

func (e *Eater) Start() {
	e.echoOut = C.cEater(C.int(e.Mem))
}

func (e *Eater) Stop() {
	C.free(unsafe.Pointer(e.echoOut))
	debug.FreeOSMemory()
	e.Running = false
}
