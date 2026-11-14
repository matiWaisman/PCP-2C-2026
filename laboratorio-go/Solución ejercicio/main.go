package main

import (
	"fmt"
	"sync"
	"time"
)

type Pago struct {
	monto        int
	nombreOrigen string
}

type Asociado struct {
	nombre            string
	yuanes            int
	pagosJefe         chan<- Pago   // Canal para enviar pagos al jefe
	comisionesJefe    <-chan int    // Canal para recibir comisiones del jefe
	pagosAsociados    <-chan Pago   // Canal para recibir pagos de los asociados
	nombresAsociados  [4]string     // Nombres de los asociados
	comisionAsociados [4]chan<- int // Canales para enviar comisiones a los asociados
}

func NuevoAsociado(nombre string, pagosJefe chan<- Pago, comisionesJefe <-chan int) *Asociado {
	return &Asociado{
		nombre:            nombre,
		yuanes:            1000,
		pagosJefe:         pagosJefe,
		comisionesJefe:    comisionesJefe,
		nombresAsociados:  [4]string{""},
		comisionAsociados: [4]chan<- int{nil, nil, nil, nil},
	}
}

func main() {
	canalNombres := make(chan string)
	go Nombres(canalNombres)
	var wg sync.WaitGroup
	wg.Add(1)
	raiz := NuevoAsociado("Hernán", nil, nil)
	go raiz.unirseGeneración(canalNombres, &wg)
	wg.Wait()
}

func (A *Asociado) nSubordinados() int {
	// Non-nil channels count as 1, nil channels count as 0
	count := 0
	for _, ch := range A.comisionAsociados {
		if ch != nil {
			count++
		}
	}
	return count
}

func (A *Asociado) indiceAsociado(nombre string) int {
	for i, n := range A.nombresAsociados {
		if n == nombre {
			return i
		}
	}
	return -1 // Not found
}

func (A *Asociado) enviarComisionesSubordinados(monto int) {
	yuanesPorAsociado := monto / A.nSubordinados()
	for _, ch := range A.comisionAsociados {
		if ch != nil {
			// ¿Qué pasa si esto no corre en una goroutine?
			go func(ch chan<- int) { ch <- yuanesPorAsociado }(ch)
		}
	}
}

func (A *Asociado) unirseGeneración(canalNombres chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	A.invitar(canalNombres, wg)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if A.nSubordinados() == 0 {
			// Si soy el jefe, imprimo mi total de yuanes y salgo
			if A.pagosJefe == nil {
				fmt.Printf("Asociado %s terminó con %d yuanes\n", A.nombre, A.yuanes)
				return
			}
			// Solo pago 100 yuanes en cada intervalo de tiempo, si no tengo asociados a los que pagarles comisiones
			select {
			case <-ticker.C:
				pago := Pago{
					monto:        100,
					nombreOrigen: A.nombre,
				}
				A.yuanes -= pago.monto
				// Si yuanes negativos, no puedo pagarle al jefe. Mando 0 y salgo
				if A.yuanes < 0 {
					pago.monto = 0
					A.pagosJefe <- pago
					return
				}
				A.pagosJefe <- pago
			case comision := <-A.comisionesJefe:
				A.yuanes += comision
			}
		} else {
			// Si tengo asociados, recibo pagos de ellos y reenvio al jefe. Salvo que sea jefe. Que me quede el 50% de la comisión y le pase el resto a los asociados
			select {
			case pago := <-A.pagosAsociados:
				// Si recibo pago 0, quiere decir que el asociado termino
				if pago.monto == 0 {
					// Encontrar el índice del asociado que me está pagando
					indice := A.indiceAsociado(pago.nombreOrigen)
					A.comisionAsociados[indice] = nil
					A.nombresAsociados[indice] = ""
				} else {
					// Reenviar el pago al jefe si tengo un jefe, sino me quedo con el 50% y le paso el resto a los asociados
					if A.pagosJefe == nil {
						A.yuanes += pago.monto / 2
						A.enviarComisionesSubordinados(pago.monto / 2)
					} else {
						A.pagosJefe <- pago
					}
				}
			case comision := <-A.comisionesJefe:
				// Si recibo comisones del jefe, me quedo con el 50% y le paso el resto a los asociados
				A.yuanes += comision / 2
				A.enviarComisionesSubordinados(comision / 2)
			}
		}
	}
}

func (A *Asociado) invitar(canalNombres chan string, wg *sync.WaitGroup) {
	// Hacer hasta 4 asociados
	pagos := make(chan Pago, 4)
	A.pagosAsociados = pagos
	for numAsociado := 0; numAsociado < 4; numAsociado++ {
		nombreAsociado, ok := <-canalNombres
		if !ok {
			return
		}
		A.nombresAsociados[numAsociado] = nombreAsociado
		fmt.Println("Asociado: ", nombreAsociado, " con jefe: ", A.nombre)
		comision := make(chan int, 1)
		A.comisionAsociados[numAsociado] = comision
		asociado := NuevoAsociado(nombreAsociado, pagos, comision)
		wg.Add(1)
		go asociado.unirseGeneración(canalNombres, wg)
	}
}
