package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Server(ch chan int, canalesClientes []chan int) {
	for i := 0; i < len(canalesClientes); i++{
		var numero int = rand.Intn(11)
		var juegoFinalizado bool = false
		var intentos int = 0
		canalesClientes[i] <- 1
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
	
}

func Cliente(chSrv chan int, miCanal chan int) {
	// Recibo un mensaje que me da el pie a ponerme a adivinar
	<- miCanal
	for i := 0; i < 11; i++ {
		chSrv <- i
		respuesta := <-chSrv
		if respuesta == 1 {
			fmt.Println("Adivine el numero ", i, " en ", i+1, " intentos")
			break
		}
	}
}

func main() {
	ch := make(chan int)
	var canales [5]chan int
	for i := range canales {
		canales[i] = make(chan int)
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		Server(ch, canales[:])
	})

	for i := 0; i < len(canales); i++{
		wg.Go(func() {
			Cliente(ch, canales[i])
		})
	}
	
	wg.Wait()
}
