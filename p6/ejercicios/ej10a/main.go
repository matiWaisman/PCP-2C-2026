package main

type Req struct {
	canal  chan int
	numero int
}

func P(canalS, canalC chan Req) {
	var cola []Req
	miCanal := make(chan int)
	var esperandoResultado, resultadoPendiente bool
	var resultado int

	for {
		var enviarS chan Req
		var mensaje Req
		if len(cola) > 0 && !esperandoResultado && !resultadoPendiente {
			enviarS = canalS
			mensaje = Req{canal: miCanal, numero: cola[0].numero}
		}

		var recibirS chan int
		if esperandoResultado {
			recibirS = miCanal
		}

		var enviarC chan int
		if resultadoPendiente {
			enviarC = cola[0].canal
		}

		select {
		case pedido := <-canalC:
			// Podemos seguir encolando mientras S trabaja o el cliente no recibe.
			cola = append(cola, pedido)
		case enviarS <- mensaje:
			esperandoResultado = true
		case resultado = <-recibirS:
			esperandoResultado = false
			resultadoPendiente = true
		case enviarC <- resultado:
			// Recién al entregar el resultado retiramos el pedido de la cola.
			cola[0] = Req{}
			cola = cola[1:]
			resultadoPendiente = false
		}
	}
}
