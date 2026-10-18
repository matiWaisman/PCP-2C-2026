package main

import (
	"container/list"
	"fmt"
)

func AsyncChan(entrada, salida chan any) {
	cola := list.New()
	var canalSalida chan any = nil
	for {
		primero := cola.Front()
		var valor any
		if primero != nil {
			valor = primero.Value
		}
		select {
		case write := <-entrada:
			cola.PushBack(write)
			canalSalida = salida
		case canalSalida <- valor:
			cola.Remove(primero)
			if cola.Len() == 0 {
				canalSalida = nil
			}
		}
	}
}

func main() {
	// Demo armada por Codex
	entrada := make(chan any)
	salida := make(chan any)
	go AsyncChan(entrada, salida)

	fmt.Printf("Capacidad de los canales: entrada=%d, salida=%d\n", cap(entrada), cap(salida))

	// Todavía no hay ningún consumidor: todos los valores quedan en la cola.
	const cantidad = 100_000
	fmt.Printf("Enviando %d valores sin que nadie lea de salida...\n", cantidad)
	for i := 0; i < cantidad; i++ {
		entrada <- i
	}
	fmt.Println("Todos los envíos terminaron antes de comenzar las lecturas.")

	// Recién ahora consumimos los valores y verificamos el orden FIFO.
	fmt.Println("Comenzando a leer...")
	for i := 0; i < cantidad; i++ {
		valor := <-salida
		if valor != i {
			panic(fmt.Sprintf("Se esperaba %d, se recibió %v", i, valor))
		}
	}
	fmt.Printf("Se recibieron los %d valores en orden, sin pérdidas.\n", cantidad)

	select {
	case valor := <-salida:
		panic(fmt.Sprintf("Se recibió un valor extra: %v", valor))
	default:
		fmt.Println("La cola está vacía: una nueva lectura quedaría bloqueada.")
	}
}
