package main

import "fmt"

// Go no permite const para slices; var a nivel de paquete es lo más cercano.
var nombres = []string{
	"Guido Ezequiel",
	"Valentin",
	"Martin Alejandro",
	"Pedro Lautaro",
	"Sofía Melina",
	"Ezequiel Ilán",
	"Juan Cruz",
	"Ivan Gonzalo",
	"Carlos Federico",
	"Ignacio Agustin",
	"Leonel",
	"Guido",
	"Mariano",
	"Ivo Martin",
	"Julián Agustín",
	"Julieta",
	"Martín",
	"Lisandro",
	"Guillermo Francisco",
	"Martín León",
	"Juan",
	"Samuel",
	"Matias",
	"Thiago Matias",
	"Tiago Martin",
	"Ezequiel David",
	"Tomas Agustin",
	"Magali",
	"Federico Uriel",
	"Mario Alejandro",
	"Mercedes",
	"Joaquín",
	"Ulises",
	"Nicolás Javier",
	"Gian Agustín",
	"Facundo Nicolás",
	"Juan Cruz",
	"Leif",
	"Damian",
	"Tobias Demián",
	"Juan Cruz",
	"Valentino",
	"Juan",
	"Santiago Ezequiel",
	"Tomas",
	"Ignacio Lucas",
	"Maura",
	"Mauricio",
	"Manuel",
	"Guilherme Heinrich",
	"Segundo",
	"Ramiro Ezequiel",
	"María Valentina",
	"Iván",
	"Ralph Ezequiel",
	"Mateo",
	"Carlos Ignacio",
	"Alexis",
	"Emiliano",
	"Bruno Emilio",
	"Pedro Nicolás",
	"Octavio Valentín",
	"Matías Eliel",
}

func Nombres(canalNombres chan string) {
	// TODO - Parte 1
	// Pista: como hay nombres repetidos, es recomendable enumerarlos para no tener problemas.
	// Go tiene un Printf casi igual al de C :)
	// fmt.Sprintf("(%d) %s", i, nombre)
	for i, nombre := range nombres {
		canalNombres <- fmt.Sprintf("(%d) %s", i, nombre)
	}
	close(canalNombres)
}
