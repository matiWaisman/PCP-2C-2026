package main

type Req struct {
	canal  chan int
	numero int
}

type Res struct {
	canal  chan int
	numero int
}

// El ticket es interno al proxy; C y S siguen intercambiando Req.
type pedidoConTicket struct {
	pedido Req
	ticket int
}

func P(canalS1, canalS2, canalC chan Req, N int) {
	var colaPendientesDeEnvio []pedidoConTicket
	var colaEnCursoS1 []pedidoConTicket
	var colaEnCursoS2 []pedidoConTicket
	// Los resultados pueden llegar en otro orden: se guardan por ticket.
	respuestasPendientes := make(map[int]Res)
	proximoTicket := 0
	siguienteAEntregar := 0
	miCanal1 := make(chan int)
	miCanal2 := make(chan int)
	var resultado int

	for {
		var enviarS1 chan Req
		var enviarS2 chan Req
		var mensajeS1, mensajeS2 Req
		if len(colaPendientesDeEnvio) > 0 && len(colaEnCursoS1) < N {
			enviarS1 = canalS1
			mensajeS1 = Req{canal: miCanal1, numero: colaPendientesDeEnvio[0].pedido.numero}
		}
		if len(colaPendientesDeEnvio) > 0 && len(colaEnCursoS2) < N {
			enviarS2 = canalS2
			mensajeS2 = Req{canal: miCanal2, numero: colaPendientesDeEnvio[0].pedido.numero}
		}

		var recibirS1 chan int
		var recibirS2 chan int
		if len(colaEnCursoS1) > 0 {
			recibirS1 = miCanal1
		}
		if len(colaEnCursoS2) > 0 {
			recibirS2 = miCanal2
		}

		var enviarC chan int
		if respuesta, listo := respuestasPendientes[siguienteAEntregar]; listo {
			resultado = respuesta.numero
			enviarC = respuesta.canal
		}

		select {
		case pedido := <-canalC:
			// Asignar el ticket al recibir conserva el orden original del cliente.
			colaPendientesDeEnvio = append(colaPendientesDeEnvio, pedidoConTicket{
				pedido: pedido,
				ticket: proximoTicket,
			})
			proximoTicket++
		case enviarS1 <- mensajeS1:
			// Podemos mandarle a S1 para que labure
			colaEnCursoS1 = append(colaEnCursoS1, colaPendientesDeEnvio[0])
			colaPendientesDeEnvio = colaPendientesDeEnvio[1:]
		case enviarS2 <- mensajeS2:
			// Podemos mandarle a S2 para que labure
			colaEnCursoS2 = append(colaEnCursoS2, colaPendientesDeEnvio[0])
			colaPendientesDeEnvio = colaPendientesDeEnvio[1:]
		case resultado = <-recibirS1:
			// S1 responde en orden: el resultado corresponde al primero de su cola.
			primero := colaEnCursoS1[0]
			respuestasPendientes[primero.ticket] = Res{canal: primero.pedido.canal, numero: resultado}
			colaEnCursoS1 = colaEnCursoS1[1:]
		case resultado = <-recibirS2:
			primero := colaEnCursoS2[0]
			respuestasPendientes[primero.ticket] = Res{canal: primero.pedido.canal, numero: resultado}
			colaEnCursoS2 = colaEnCursoS2[1:]
		case enviarC <- resultado:
			delete(respuestasPendientes, siguienteAEntregar)
			siguienteAEntregar++
		}
	}
}
