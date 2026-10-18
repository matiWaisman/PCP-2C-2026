package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Mensaje struct {
	intento int
	canal   chan int
}

func inicializarNumeros(longitud int) []int {
	res := make([]int, longitud)
	for i := range longitud {
		res[i] = rand.Intn(11)
	}
	return res
}

func Server(canalesClientes [3]chan Mensaje) {
	cantidadIntentos := make([]int, 3)
	numeros := inicializarNumeros(3)
	cantidadAdivinados := 0

	procesarIntento := func(i int, intento Mensaje) {
		cantidadIntentos[i]++
		if intento.intento == numeros[i] {
			intento.canal <- cantidadIntentos[i]
			fmt.Println(i, " Adivino el numero ", intento.intento, " en ", cantidadIntentos[i], " intentos")
			cantidadAdivinados++
			// Un canal nil deshabilita su caso en el select.
			canalesClientes[i] = nil
		} else {
			intento.canal <- -1
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

func Cliente(canalServer chan Mensaje) {
	canal := make(chan int, 11)
	// Una permutación contiene todos los números del 0 al 10 sin repetir.
	numeros := rand.Perm(11)
	siguiente := 0
	for {
		var canalEnvio chan Mensaje
		var mensaje Mensaje
		if siguiente < len(numeros) {
			canalEnvio = canalServer
			mensaje = Mensaje{intento: numeros[siguiente], canal: canal}
		}
		// Al agotar los intentos, canalEnvio queda nil y solo recibimos respuestas.
		select {
		case canalEnvio <- mensaje:
			siguiente++
		case respuesta := <-canal:
			if respuesta != -1 {
				return
			}
		}
	}
}

func main() {
	var canales [3]chan Mensaje
	for i := range canales {
		canales[i] = make(chan Mensaje, 11)
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		Server(canales)
	})

	for i := 0; i < len(canales); i++ {
		wg.Go(func() {
			Cliente(canales[i])
		})
	}

	wg.Wait()
}
