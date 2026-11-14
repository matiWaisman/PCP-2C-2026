/*
Ejercicio 7 b)

La solución sigue siendo correcta con comunicación asincrónica porque todos los
pedidos se procesan secuencialmente desde un mismo canal FIFO. Cada lectura
devuelve el valor de la última escritura procesada antes de ella y el cliente
espera esa respuesta mediante un canal propio.

Terminar de enviar una escritura no significa que el servidor ya la haya
ejecutado, pero una lectura enviada después por el mismo canal será procesada
después de esa escritura. Otros clientes pueden intercalar escrituras.
Si se quisiera garantizar que la escritura fue ejecutada antes de que el cliente
continúe, sería necesario agregar una confirmación del servidor.
*/
package main

import (
	"fmt"
	"sync"
)

type Accion string

const (
	Escribir Accion = "Escribir"
	Leer     Accion = "Leer"
)

type Pedido struct {
	accion Accion
	valor  any
	canal  chan any
}

func Server(ch chan Pedido) {
	var compartida any
	for {
		pedido := <-ch
		if pedido.accion == Leer {
			pedido.canal <- compartida
		}
		if pedido.accion == Escribir {
			compartida = pedido.valor
		}
	}
}

func leer(ch chan Pedido) any {
	respuesta := make(chan any)
	ch <- Pedido{accion: Leer, canal: respuesta}
	return <-respuesta
}

func Cliente(id int, ch chan Pedido) {
	valor := id * 10
	ch <- Pedido{accion: Escribir, valor: valor}
	// Otro cliente puede escribir entre nuestra escritura y nuestra lectura.
	leido := leer(ch)
	fmt.Printf("Cliente %d: escribió %d y después leyó %v\n", id, valor, leido)
}

func main() {
	pedidos := make(chan Pedido)
	go Server(pedidos)

	fmt.Println("Valor inicial:", leer(pedidos))
	pedidos <- Pedido{accion: Escribir, valor: 42}
	fmt.Println("Después de escribir 42:", leer(pedidos))

	fmt.Println("Tres clientes acceden concurrentemente a la misma variable.")
	fmt.Println("Cada lectura devuelve el valor de la última escritura procesada por el servidor.")
	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		wg.Go(func() {
			Cliente(id, pedidos)
		})
	}
	wg.Wait()

	fmt.Println("Valor final:", leer(pedidos))
}
