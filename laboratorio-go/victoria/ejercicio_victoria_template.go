package main

import (
	"fmt"
	"math/rand"
	"time"
)
// En este caso usamos el mismo struct para representar pedidos de faros y carrocerias porque no hace
// falta diferenciarlos a fines de este ejercicio.
type Pedido struct {
	colaboradorID int // Se utiliza para el print en el método Victoria
	terminado     chan struct{}
}

type Taller struct {
	faros       chan Pedido
	carrocerias chan Pedido
}

func NuevoTaller() *Taller {
	canalFaros := make(chan Pedido)
	canalCarrocerias := make(chan Pedido)
	return &Taller{faros: canalFaros, carrocerias: canalCarrocerias}
}

func (t *Taller) TraerFaro(colaboradorID int) {
	miCanal := make(chan struct{})
	for{
		miPedido := Pedido{colaboradorID: colaboradorID, terminado: miCanal}
		t.faros <- miPedido
		<- miCanal 
	}
}

func (t *Taller) TraerCarroceria(colaboradorID int) {
	miCanal := make(chan struct{})
	for{
		miPedido := Pedido{colaboradorID: colaboradorID, terminado: miCanal}
		t.carrocerias <- miPedido
		<- miCanal 
	}
}

func armarUnidad() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func (t *Taller) Victoria(cantidadAutos int) {
	for range cantidadAutos{
		var cantidadFaros int = 0 
		var cantidadCarroceria int = 0
		canalCarrocerias := t.carrocerias
		canalFaros := t.faros
		var idCarroceria, idFaro1, idFaro2 int
		var canalAProductorCarroceria chan(struct{}) = nil 
		var canalAProductorFaro1 chan(struct{}) = nil 
		var canalAProductorFaro2 chan(struct{}) = nil 
		for (cantidadFaros < 2 || cantidadCarroceria < 1){
			select{
				case pedidoCarroceria := <- canalCarrocerias: 
					cantidadCarroceria += 1
					canalCarrocerias = nil
					canalAProductorCarroceria = pedidoCarroceria.terminado
					idCarroceria = pedidoCarroceria.colaboradorID
				case pedidoFaro := <- canalFaros: 
					cantidadFaros += 1
					if cantidadFaros == 1{
						canalAProductorFaro1 = pedidoFaro.terminado
						idFaro1 = pedidoFaro.colaboradorID
					} else{
						canalAProductorFaro2 = pedidoFaro.terminado
						idFaro2 = pedidoFaro.colaboradorID
						canalFaros = nil
					}
			}
		}
		armarUnidad()
		canalAProductorCarroceria <- struct{}{}
		canalAProductorFaro1 <- struct{}{}
		canalAProductorFaro2 <- struct{}{}
		fmt.Printf("Auto armado: carrocería del productor %d, faros de los productores %d y %d\n", idCarroceria, idFaro1, idFaro2)
	}
}
