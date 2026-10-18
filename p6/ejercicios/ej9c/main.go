package main

type Pedido string

const (
	String Pedido = "String"
	Canal  Pedido = "Canal"
)

type Mensaje struct {
	pedido Pedido
	cadena string
	canal  chan string
}

// Cifrado César: calcula el desplazamiento a partir de la clave k.
func codificar(m, k string) string {
	var desplazamiento rune
	for _, caracter := range k {
		desplazamiento = (desplazamiento + caracter) % 26
	}

	resultado := []rune(m)
	for i, letra := range resultado {
		switch {
		case letra >= 'a' && letra <= 'z':
			resultado[i] = 'a' + (letra-'a'+desplazamiento)%26
		case letra >= 'A' && letra <= 'Z':
			resultado[i] = 'A' + (letra-'A'+desplazamiento)%26
		}
	}
	return string(resultado)
}

// No se anticipa a las claves k, las pide recien cuando tiene un mensaje
func T(chCliente chan Mensaje, chServicio chan chan string) {
	var canalAEnviar chan string
	canalARecibirK := make(chan string)
	for {
		mensaje := <-chCliente
		if mensaje.pedido == String {
			chServicio <- canalARecibirK
			k := <-canalARecibirK
			canalAEnviar <- codificar(mensaje.cadena, k)
		}
		if mensaje.pedido == Canal {
			canalAEnviar = mensaje.canal
		}
	}
}
