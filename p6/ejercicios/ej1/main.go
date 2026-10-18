package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Cmd string

const (
	Next  Cmd = "Next"
	Count Cmd = "Count"
)

func Srv(ch chan Cmd) {
	var count int = 0
	for {
		pedido := <-ch
		switch pedido {
		case Next:
			count = count + 1
		case Count:
			fmt.Println("La cantidad de counts fue: ", count)
			count = 0
		}
	}
}

func Client(n int, ch chan Cmd) {
	for {
		for i := 0; i < n; i++ {
			ch <- Next
		}
		ch <- Count
	}
}

func main() {
	canal := make(chan Cmd)
	var wg sync.WaitGroup

	wg.Go(func() {
		Srv(canal)
	})

	wg.Go(func() {
		Client(rand.Intn(10)+1, canal)
	})

	wg.Go(func() {
		Client(rand.Intn(10)+1, canal)
	})

	wg.Wait()
}
