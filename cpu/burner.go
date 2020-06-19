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

// Burner is a CPU burner that runs a goroutine in a core with the TTL specified
type Burner struct {
	Running bool      `json:"running"`
	NumBurn int       `json:"num_burn"`
	TTL     int       `json:"ttl"`
	UUID    uuid.UUID `json:"id,omitempty"`
	chans   []chan bool
}

// NewBurner returns a new Burner
func NewBurner(body io.ReadCloser) (handler.Burner, error) {
	c := Burner{}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&c)
	c.UUID = uuid.NewV4()
	c.Running = true

	c.chans = make([]chan bool, c.NumBurn)
	for i := range c.chans {
		c.chans[i] = make(chan bool)
	}

	return &c, err
}

// ID returns the ID of the Burner
func (c Burner) ID() uuid.UUID {
	return c.UUID
}

// IsRunning checks if the TTL is over or if the Burner has been stopped
func (c Burner) IsRunning() bool {
	start := time.Now()
	isDone := time.Since(start) < time.Millisecond*time.Duration(c.TTL)
	return isDone && c.Running
}

// Start runs the given number of goroutines to burn a the specified CPUs
func (c *Burner) Start() {
	fmt.Printf("Burning %d CPUs/cow\n", c.NumBurn)

	for i := 0; i < c.NumBurn; i++ {
		go cpuBurn(c.chans[i], i)
		c.chans[i] <- true
	}
	fmt.Printf("Sleeping %d miliseconds\n", c.TTL)
	for c.IsRunning() {
		for i := range c.chans {
			c.chans[i] <- true

		}
	}
	c.Stop()
}

// Stop stops the Burner
func (c *Burner) Stop() {
	c.Running = false

	for i := range c.chans {
		c.chans[i] <- false
	}

	c.NumBurn = 0
	c.TTL = 0
}

func cpuBurn(cont chan bool, i int) {
	for <-cont {
		for i := 0; i < 2147483647; i++ {
		}
		runtime.Gosched()
	}
}
