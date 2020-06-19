package cpu

import (
	"encoding/json"

	"github.com/JavierClairvaux/px-usewithcare/handler"
	uuid "github.com/satori/go.uuid"

	"fmt"
	"io"
	"runtime"
	"time"
)

type Burner struct {
	Running bool      `json:"running"`
	NumBurn int       `json:"num_burn"`
	TTL     int       `json:"ttl"`
	UUID    uuid.UUID `json:"id,omitempty"`
}

func NewBurner(body io.ReadCloser) (handler.Burner, error) {
	c := Burner{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&c)
	c.UUID = uuid.NewV4()
	c.Running = true
	return &c, err
}

func (c Burner) ID() uuid.UUID {
	return c.UUID
}

func (c Burner) IsRunning() bool {
	return c.Running
}

func (c *Burner) Start() {
	fmt.Printf("Burning %d CPUs/cow\n", c.NumBurn)
	for i := 0; i < c.NumBurn; i++ {
		go cpuBurn(c)
	}
	fmt.Printf("Sleeping %d miliseconds\n", c.TTL)
	for start := time.Now(); time.Since(start) < time.Millisecond*time.Duration(c.TTL); {
	}
	c.NumBurn = 0
	c.Running = false
	c.TTL = 0
}

func (c *Burner) Stop() {
	c.NumBurn = 0
	c.Running = false
	c.TTL = 0
}

func cpuBurn(c *Burner) {
	for {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
		if !c.Running {
			break
		}
	}
}
