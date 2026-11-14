package main

import (
	"container/list"
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
	var canalFaros chan Pedido = make(chan Pedido)
	var canalCarrocerias chan Pedido = make(chan Pedido)
	return &Taller{
		faros:       canalFaros,
		carrocerias: canalCarrocerias,
	}
}

func (t *Taller) TraerFaro(colaboradorID int) {
	var canalRespuesta chan struct {} = make(chan struct{})
	pedido := Pedido{
		colaboradorID: colaboradorID,
		terminado: canalRespuesta,
	}
	t.faros <- pedido
	<- canalRespuesta
}

func (t *Taller) TraerCarroceria(colaboradorID int) {
	var canalRespuesta chan struct {} = make(chan struct{})
	pedido := Pedido{
		colaboradorID: colaboradorID,
		terminado: canalRespuesta,
	}
	t.carrocerias <- pedido
	<- canalRespuesta
}

func armarUnidad() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func (t *Taller) Victoria(cantidadAutos int) {
	canalFaros := t.faros
	canalCarrocerias := t.carrocerias
	cantidadFarosRecibidos := 0
	colaCarrocerias := list.New()
	colaFaros := list.New()
	colaADespertar := list.New()
	autosArmados := 0
	for autosArmados < cantidadAutos{
		if cantidadFarosRecibidos == 2 && canalCarrocerias == nil{
			// Tenemos los materiales para laburar 
			armarUnidad()
			autosArmados += 1
			cantidadFarosRecibidos = 0
			canalFaros = t.faros
			canalCarrocerias = t.carrocerias
			for i := 0; i < 2; i++ {
				colaADespertar.PushBack(colaFaros.Remove(colaFaros.Front()))
			}
			colaADespertar.PushBack(colaCarrocerias.Remove(colaCarrocerias.Front()))
			if autosArmados == cantidadAutos {
				break
			}
		}
		aDespertar := colaADespertar.Front()
		var canalDespertar chan struct{}
		if aDespertar != nil {
			canalDespertar = aDespertar.Value.(chan struct{})
		}

		select{
		case pedido := <- canalFaros:
			cantidadFarosRecibidos += 1 
			if cantidadFarosRecibidos == 2 {
				canalFaros = nil 
			}
			colaFaros.PushBack(pedido.terminado)
			fmt.Printf("Faro recibido de colaborador %d\n", pedido.colaboradorID)
		
		case pedido := <- canalCarrocerias:
			canalCarrocerias = nil 
			colaCarrocerias.PushBack(pedido.terminado)
			fmt.Printf("Carroceria recibido de colaborador %d\n", pedido.colaboradorID)
			
		case canalDespertar <- struct{}{}:
			colaADespertar.Remove(aDespertar)	
		}
	}
	// Despierto a los ultimos 3 manualmente 
	for colaADespertar.Len() > 0 {
		canalDespertar := colaADespertar.Remove(colaADespertar.Front()).(chan struct{})
		canalDespertar <- struct{}{}
	}
}
