package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func inicializarNumeros(longitud int) []int {
	res := make([]int, longitud)
	for i := 0; i < longitud; i++ {
		res[i] = rand.Intn(11)
	}
	return res
}

func ServerBloqueante(canalesClientes [3]chan int) {
	cantidadIntentos := make([]int, 3)
	numeros := inicializarNumeros(3)
	cantidadAdivinados := 0

	procesarIntento := func(i, intento int) {
		cantidadIntentos[i]++
		if intento == numeros[i] {
			canalesClientes[i] <- 1
			fmt.Println(i, " Adivino el numero ", intento, " en ", cantidadIntentos[i], " intentos")
			cantidadAdivinados++
			// Un canal nil deshabilita su caso en el select.
			canalesClientes[i] = nil
		} else {
			canalesClientes[i] <- 0
		}
	}

	for cantidadAdivinados < len(canalesClientes) {
		select {
		case intento := <-canalesClientes[0]:
			procesarIntento(0, intento)
		case intento := <-canalesClientes[1]:
			procesarIntento(1, intento)
		case intento := <-canalesClientes[2]:
			procesarIntento(2, intento)
		}
	}
}

// Espera activa, si estan todos los canales vacios volves a consultar a los canales
func Server(canalesClientes []chan int) {
	cantidadClientes := len(canalesClientes)
	cantidadIntentos := make([]int, cantidadClientes)
	adivino := make([]bool, cantidadClientes)
	numeros := inicializarNumeros(cantidadClientes)
	cantidadAdivinados := 0
	for cantidadAdivinados < cantidadClientes {
		for i := 0; i < len(canalesClientes); i++ {
			if !adivino[i] {
				select {
				case intento := <-canalesClientes[i]:
					cantidadIntentos[i] += 1
					if intento == numeros[i] {
						adivino[i] = true
						cantidadAdivinados++
						canalesClientes[i] <- 1
						fmt.Println(i, " Adivino el numero ", intento, " en ", cantidadIntentos[i], " intentos")
					} else {
						canalesClientes[i] <- 0
					}
				default:
					// No hay un mensaje disponible: seguir con el próximo.
				}
			}
		}
	}
}

func Cliente(miCanal chan int) {
	for i := 0; i < 11; i++ {
		miCanal <- i
		respuesta := <-miCanal
		if respuesta == 1 {
			fmt.Println("Adivine el numero ", i, " en ", i+1, " intentos")
			break
		}
	}
}

func main() {
	var canales [3]chan int
	for i := range canales {
		canales[i] = make(chan int)
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		ServerBloqueante(canales)
	})

	for i := 0; i < len(canales); i++ {
		wg.Go(func() {
			Cliente(canales[i])
		})
	}

	wg.Wait()
}
