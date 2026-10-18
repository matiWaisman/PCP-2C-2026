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

type Pedido struct {
	cmd   Cmd
	canal chan int
}

func Srv(ch chan Pedido) {
	var count int = 0
	for {
		var pedido Pedido = <-ch
		switch pedido.cmd {
		case Next:
			count = count + 1
		case Count:
			pedido.canal <- count
			count = 0
		}
	}
}

func Client(n int, ch chan Pedido) {
	miCanal := make(chan int)
	for {
		for i := 0; i < n; i++ {
			ch <- Pedido{
				cmd: Next,
			}
		}
		ch <- Pedido{
			cmd:   Count,
			canal: miCanal,
		}
		count := <-miCanal
		fmt.Println("La cantidad de mensajes Next fue:", count)
	}
}

func main() {
	canal := make(chan Pedido)
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
