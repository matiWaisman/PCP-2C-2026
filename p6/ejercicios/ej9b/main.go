package main

type Pedido string

const (
	String Pedido = "String"
	Canal Pedido = "Canal"
)

type Mensaje struct{
	pedido 		Pedido
	cadena 		string
	canal 		chan string
}

// Cifrado César: desplaza tres posiciones las letras A-Z y a-z.
func codificar(m string) string {
	resultado := []rune(m)
	for i, letra := range resultado {
		switch {
		case letra >= 'a' && letra <= 'z':
			resultado[i] = 'a' + (letra-'a'+3)%26
		case letra >= 'A' && letra <= 'Z':
			resultado[i] = 'A' + (letra-'A'+3)%26
		}
	}
	return string(resultado)
}

func T(ch chan Mensaje){
	var canalAEnviar chan string = nil
	for{
		mensaje := <- ch 
		if mensaje.pedido == String{
			canalAEnviar <- codificar(mensaje.cadena)
		}
		if mensaje.pedido == Canal{
			canalAEnviar = mensaje.canal
		}

	}
}