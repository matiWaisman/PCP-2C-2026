package main

import "time"

// Opcion si no me preocupo si quedan mensajes pendientes que no se pudieron enviar x tener el buffer lleno
func Timer(frecuencia int, canales []chan struct{}) {
	cantidadCanales := len(canales)
	ticker := time.NewTicker(time.Duration(frecuencia) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		for i := range cantidadCanales {
			select {
			case canales[i] <- struct{}{}:
				// Se pudo enviar
			default:
				// No hacemos nada
			}
		}
	}
}

// Opcion si me preocupo por todos los mensajes que no se pudieron enviar
func Timer2(frecuencia int, canales []chan struct{}) {
	cantidadCanales := len(canales)
	pendientes := make([]int, cantidadCanales)
	ticker := time.NewTicker(time.Duration(frecuencia) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		for i := range cantidadCanales {
			pendientes[i] += 1
			sigo := true
			for pendientes[i] > 0 && sigo {
				select {
				case canales[i] <- struct{}{}:
					// Se pudo enviar
					pendientes[i] -= 1
				default:
					// El buffer esta lleno, asi que no perdemos tiempo y seguimos con el resto
					sigo = false
				}
			}
		}
	}
}
