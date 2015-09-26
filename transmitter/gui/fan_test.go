package gui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFanIn(t *testing.T) {
	assert := assert.New(t)
	expected := make([]int, 1000)
	for i := 0; i < len(expected); i++ {
		expected[i] = i
	}

	channels := make([]chan int, 13)
	outputs := make([]chan<- int, len(channels))
	inputs := make([]<-chan int, len(channels))
	for i := range channels {
		channels[i] = make(chan int)
		outputs[i] = channels[i]
		inputs[i] = channels[i]
	}

	out := FanIn(inputs...)
	go func() {
		for i := 0; i < len(expected); i++ {
			outputs[i%len(outputs)] <- expected[i]
		}
		for _, c := range outputs {
			close(c)
		}
	}()
	counter := 0
	for c := range out {
		assert.Contains(expected, c)
		counter++
	}
	assert.Equal(len(expected), counter)
}
