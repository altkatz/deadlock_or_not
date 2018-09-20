// +build race
package main

import (
	"log"
	_ "net/http"
	"time"
	// "strconv"
)

type RingBuffer struct {
	inputChannel  <-chan int
	outputChannel chan int
}

func NewRingBuffer(inputChannel <-chan int, outputChannel chan int) *RingBuffer {
	return &RingBuffer{inputChannel, outputChannel}
}

func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		select {
		case r.outputChannel <- v:
			log.Printf("send %d to out chan\n", v)
		default:
			old := <-r.outputChannel
			r.outputChannel <- v
			log.Printf("discard %d and send %d \n", old, v)
		}
	}
	close(r.outputChannel)
}

type server struct {
	inchan  chan int
	outchan chan int
}

func main() {
	in := make(chan int)
	out := make(chan int, 5)
	time.Sleep(1 * time.Second)
	rb := NewRingBuffer(in, out)
	go rb.Run()
	for i := 0; i < 10; i++ {
		go func(i int) {
			in <- i
		}(i)
	}
	for res := range out {
		log.Printf("from_out_put_chan:%d\n", res)
	}
}
