package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Server(ch chan int) {
	var numero int = rand.Intn(11)
	var juegoFinalizado bool = false
	var intentos int = 0
	for !juegoFinalizado {
		intentos += 1
		var intento int = <-ch
		if intento == numero {
			juegoFinalizado = true
			fmt.Println("Adivino el numero ", numero, " en ", intentos, " intentos")
			ch <- 1
		} else {
			ch <- 0
		}
	}
}

func Cliente(ch chan int) {
	for i := 0; i < 11; i++ {
		ch <- i
		respuesta := <-ch
		if respuesta == 1 {
			fmt.Println("Adivine el numero ", i, " en ", i+1, " intentos")
			break
		}
	}
}

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup

	wg.Go(func() {
		Server(ch)
	})

	wg.Go(func() {
		Cliente(ch)
	})

	wg.Wait()
}
