package gui

import (
	"fmt"
	"time"
	"github.com/pivotal-golang/bytefmt"
)

func renderProgress(value float64) {
	fmt.Printf("%f%%", value)
}

func renderSpeed(value uint64) {
	fmt.Println(" ", bytefmt.ByteSize(value), "/sec")
}

type Gui struct {
	Total uint64
}

func (this *Gui) Loop(input <-chan int) {
	var accumulator uint64 = 0
	beginning := time.Now()
	var oldDelta int64
	for n := range input {
		delta := time.Since(beginning).Nanoseconds()
		accumulator += uint64(n)
		if (delta - oldDelta) >= 250000000 {
			speed := accumulator * 1000000000 / uint64(delta)
			renderProgress(float64(accumulator * 10000 / this.Total)/100)
			renderSpeed(speed)
			oldDelta = delta
		}
	}
}
