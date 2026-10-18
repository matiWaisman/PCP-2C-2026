package main

type Req struct {
	canal  chan int
	numero int
}

type Res struct {
	canal 	chan int 
	numero	int 
}

func P(canalS, canalC chan Req, N int) {
	var colaPendientesDeEnvio []Req 
	var colaEnCurso []Req 
	var respuestasPendientes []Res 
	miCanal := make(chan int)
	var resultado int

	for {
		var enviarS chan Req
		var mensaje Req
		if len(colaPendientesDeEnvio) > 0 && len(colaEnCurso) < N {
			enviarS = canalS
			mensaje = Req{canal: miCanal, numero: colaPendientesDeEnvio[0].numero}
		}

		var recibirS chan int
		if len(colaEnCurso) > 0 {
			recibirS = miCanal
		}

		var enviarC chan int
		if len(respuestasPendientes) > 0 {
			resultado = respuestasPendientes[0].numero
			enviarC = respuestasPendientes[0].canal
		}

		select {
		case pedido := <-canalC:
			// Recibimos un pedido que encolamos
			colaPendientesDeEnvio = append(colaPendientesDeEnvio, pedido)
		case enviarS <- mensaje:
			// Podemos mandarle a S para que labure
			colaEnCurso = append(colaEnCurso, colaPendientesDeEnvio[0])
			colaPendientesDeEnvio = colaPendientesDeEnvio[1:]
		case resultado = <-recibirS:
			// Recibimos un resultado de S
			respuestasPendientes = append(respuestasPendientes, Res{canal: colaEnCurso[0].canal, numero: resultado})
			colaEnCurso = colaEnCurso[1:]
		case enviarC <- resultado:
			// Le mandamos a C su resultado
			respuestasPendientes = respuestasPendientes[1:]
		}
	}
}
