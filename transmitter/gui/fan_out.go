package gui

import ()

func FanOut(input <-chan int, count int) []<-chan int {
	channels := make([]chan int, count)
	outputs := make([]<-chan int, count)
	for i := range channels {
		channels[i] = make(chan int)
		outputs[i] = channels[i]
	}
	go func() {
		for n := range input {
			for _, c := range channels {
				c <- n
			}
		}
		for _, c := range channels {
			close(c)
		}
	}()
	return outputs
}
