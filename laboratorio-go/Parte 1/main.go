package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Pago struct {
	monto        int
	nombreOrigen string
}

type Asociado struct {
	nombre              string
	yuanes              int
	pagosJefe           chan<- Pago
	comisionesJefe      <-chan int
	canalesSubordinados [4]chan<- int
	nombresSubordinados [4]string
}

type PagoComision struct {
	monto int
	canal chan<- int
}

func NuevoAsociado(nombre string, pagosJefe chan<- Pago, comisionesJefe <-chan int) *Asociado {
	return &Asociado{
		nombre:              nombre,
		yuanes:              1000,
		pagosJefe:           pagosJefe,
		comisionesJefe:      comisionesJefe,
		canalesSubordinados: [4]chan<- int{},
		nombresSubordinados: [4]string{},
	}
}

func main() {
	hernan := NuevoAsociado("Hernan", nil, nil)
	var wg sync.WaitGroup
	var canalNombres chan string = make(chan string)
	wg.Go(func() {
		Nombres(canalNombres)
	})

	wg.Go(func() {
		hernan.unirseGeneración(canalNombres, &wg)
	})

	wg.Wait()
}

func (A *Asociado) unirseGeneración(canalNombres chan string, wg *sync.WaitGroup) {
	canalPagos := make(chan Pago)
	A.invitar(canalNombres, canalPagos, wg)
	A.recaudarConSubordinados(canalPagos)

	if A.pagosJefe == nil {
		fmt.Printf("Soy %s y el total es %d yuanes\n", A.nombre, A.yuanes)
		return
	}
	A.pagarSinSubordinados()
}

func (A *Asociado) recaudarConSubordinados(canalPagos <-chan Pago) {
	colaComisiones := list.New()
	for {
		A.eliminarComisionesDeRetirados(colaComisiones)
		if A.nSubordinados() == 0 {
			return
		}

		// Si la cola está vacía, el canal nil deshabilita el envío del select.
		pendiente := PagoComision{}
		primerNodo := colaComisiones.Front()
		if primerNodo != nil {
			pendiente = primerNodo.Value.(PagoComision)
		}
		select {
		case pago := <-canalPagos:
			if pago.monto == 0 {
				A.anularCanal(pago.nombreOrigen)
			} else if A.pagosJefe != nil {
				A.pagosJefe <- pago
			} else {
				A.yuanes += pago.monto
				A.encolarComisiones(pago.monto/2, colaComisiones)
			}
		case comision := <-A.comisionesJefe:
			A.yuanes += comision
			A.encolarComisiones(comision/2, colaComisiones)
		case pendiente.canal <- pendiente.monto:
			A.yuanes -= pendiente.monto
			colaComisiones.Remove(primerNodo)
		}
	}
}

func (A *Asociado) encolarComisiones(monto int, cola *list.List) {
	cantidad := A.nSubordinados()
	if cantidad == 0 {
		return
	}
	for _, canal := range A.canalesSubordinados {
		if canal != nil {
			cola.PushBack(PagoComision{monto: monto / cantidad, canal: canal})
		}
	}
}

func (A *Asociado) eliminarComisionesDeRetirados(cola *list.List) {
	for nodo := cola.Front(); nodo != nil; {
		siguiente := nodo.Next()
		pendiente := nodo.Value.(PagoComision)
		activo := false
		for _, canal := range A.canalesSubordinados {
			if canal != nil && canal == pendiente.canal {
				activo = true
				break
			}
		}
		if !activo {
			cola.Remove(nodo)
		}
		nodo = siguiente
	}
}

func (A *Asociado) pagarSinSubordinados() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case comision := <-A.comisionesJefe:
			A.yuanes += comision
		case <-ticker.C:
			if A.yuanes < 100 {
				A.pagosJefe <- Pago{monto: 0, nombreOrigen: A.nombre}
				return
			}
			A.pagosJefe <- Pago{monto: 100, nombreOrigen: A.nombre}
			A.yuanes -= 100
		}
	}
}

func (A *Asociado) anularCanal(nombre string) {
	for i := range 4 {
		if A.nombresSubordinados[i] == nombre {
			A.canalesSubordinados[i] = nil
			break
		}
	}
}

func (A *Asociado) nSubordinados() int {
	res := 0
	for i := range 4 {
		if A.canalesSubordinados[i] != nil {
			res += 1
		}
	}
	return res
}

func (A *Asociado) invitar(canalNombres chan string, pagosAsociados chan Pago, wg *sync.WaitGroup) {
	for i := range 4 {
		var misComisiones chan int = make(chan int)
		nombreActual, ok := <-canalNombres
		if !ok {
			break
		}
		asociado := NuevoAsociado(nombreActual, pagosAsociados, misComisiones)
		A.canalesSubordinados[i] = misComisiones
		A.nombresSubordinados[i] = nombreActual
		wg.Go(func() {
			asociado.unirseGeneración(canalNombres, wg)
		})
	}
}
