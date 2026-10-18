package main

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

func T(entrada, salida chan string){
	for{
		mensaje := <- entrada 
		salida <- codificar(mensaje)
	}
}