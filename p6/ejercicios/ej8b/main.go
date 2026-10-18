package main

import (
	"fmt"
	"time"
)

// Opcion si no me preocupo si quedan mensajes pendientes que no se pudieron enviar x tener el buffer lleno
func Timer(frecuencia int, registros chan chan struct{}) {
	var canales []chan struct{}
	ticker := time.NewTicker(time.Duration(frecuencia) * time.Millisecond)
	defer ticker.Stop()
	for {
		select{
		case canal := <- registros:
			canales = append(canales,canal)
		case <- ticker.C: 
			for i := range(canales){
				select{
				case canales[i] <- struct{}{}:
					// Se pudo enviar
				default:
					// No se envia nada
				}
			}
		}
	}
}



func Cliente(canalSrv chan chan struct{}) {
	miCanal := make(chan struct{}, 3)
	canalSrv <- miCanal
	registro := time.Now()
	fmt.Printf("Cliente: registro enviado a las %s\n", registro.Format("15:04:05.000"))
	for {
		<-miCanal
		recepcion := time.Now()
		impresion := time.Now()
		fmt.Printf("Tick recibido a las %s | desde el registro: %s | print a las %s\n",
			recepcion.Format("15:04:05.000"),
			recepcion.Sub(registro).Round(time.Millisecond),
			impresion.Format("15:04:05.000"))
	}
}
