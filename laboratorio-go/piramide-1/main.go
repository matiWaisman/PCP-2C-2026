package main

import "fmt"

type Pago struct {
	monto        int
	nombreOrigen string
}

type Asociado struct {
	nombre            string
	yuanes            int
}

func NuevoAsociado(nombre string, pagosJefe chan<- Pago, comisionesJefe <-chan int) *Asociado {
	return &Asociado{
		nombre:            nombre,
		yuanes:            1000,
	}
}

func main() {
	// TODO
}


func (A *Asociado) invitar(canalNombres chan string) {
	miCanal := make(chan Asociado)
	
}
