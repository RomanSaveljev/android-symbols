package gui

import (
	"sync"
)

func FanIn(inputs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(inputs))

	for _, c := range inputs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
