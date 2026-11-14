package main

import "container/list"

type Req struct {
	entero         int
	canalRespuesta chan int
}

/*
func p(canalC chan Req, canalS chan Req){
	cola := list.New()
	var canalRecepcion chan int = make(chan int)
	var canalClienteEsperando chan int = nil
	var canalAEnviarleClienteEsperando chan int = nil
	var respuestaGuardada int
	for{
		primero := cola.Front()
		// Si no estamos esperando ningun pedido, es pq no se le envio nada al server, asi que le mandamos
		if canalClienteEsperando == nil && primero != nil {
			pedido := primero.Value.(Req)
			canalClienteEsperando = pedido.canalRespuesta

			req := Req{
				entero:         pedido.entero,
				canalRespuesta: canalRecepcion,
			}
			canalS <- req
		}
		select{
			case pedido := <- canalC:
				cola.PushBack(pedido)
			case respuesta := <- canalRecepcion:
				respuestaGuardada = respuesta
				canalAEnviarleClienteEsperando = canalClienteEsperando
			case canalAEnviarleClienteEsperando <- respuestaGuardada:
				cola.Remove(primero)
				canalClienteEsperando = nil
				canalAEnviarleClienteEsperando = nil
		}
	}
}
*/

type Respuesta struct {
	resultado    int
	canalCliente chan int
}

func p(canalC chan Req, canalS chan Req) {
	colaPendientesEnvio := list.New()
	colaRespuestas := list.New()
	canalRecepcion := make(chan int)
	var ocupado bool
	var canalClienteEnCurso chan int

	for {
		primeroPendienteEnvio := colaPendientesEnvio.Front()
		var canalEnvio chan Req
		var pedidoAEnviar Req
		// Solo habilitamos el envio si el servidor esta libre y hay pedidos.
		if !ocupado && primeroPendienteEnvio != nil {
			pedido := primeroPendienteEnvio.Value.(Req)
			canalEnvio = canalS
			pedidoAEnviar = Req{
				entero:         pedido.entero,
				canalRespuesta: canalRecepcion,
			}
		}

		var canalRespuestaServidor chan int
		if ocupado {
			canalRespuestaServidor = canalRecepcion
		}

		primeraRespuesta := colaRespuestas.Front()
		var canalEntrega chan int
		var resultadoAEntregar int
		// Si no hay respuestas pendientes, canalEntrega queda nil.
		if primeraRespuesta != nil {
			respuesta := primeraRespuesta.Value.(Respuesta)
			canalEntrega = respuesta.canalCliente
			resultadoAEntregar = respuesta.resultado
		}

		select {
		case pedido := <-canalC:
			colaPendientesEnvio.PushBack(pedido)

		case canalEnvio <- pedidoAEnviar:
			pedido := primeroPendienteEnvio.Value.(Req)
			canalClienteEnCurso = pedido.canalRespuesta
			colaPendientesEnvio.Remove(primeroPendienteEnvio)
			ocupado = true

		case resultado := <-canalRespuestaServidor:
			colaRespuestas.PushBack(Respuesta{
				resultado:    resultado,
				canalCliente: canalClienteEnCurso,
			})
			// El servidor puede recibir otro pedido aunque el cliente no lea aun.
			ocupado = false
			canalClienteEnCurso = nil

		case canalEntrega <- resultadoAEntregar:
			colaRespuestas.Remove(primeraRespuesta)
		}
	}
}
