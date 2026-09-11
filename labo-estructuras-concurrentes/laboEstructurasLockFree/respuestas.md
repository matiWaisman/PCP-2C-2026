# Ejercicio 3 
Lo importante es que todos sigan el orden. 
Mientras todos vayan para un lado o para el otro todo bien. 
O vas de izquierda a derecha agarrando predecesor y actual 
O vas de derecha a izquierda tomando posterior y predecesor. 

El orden es tomar primero la cabeza, el siguiente, liberar la cabeza y tomar el siguiente del siguiente y liberar el siguiente. 

En el metodo contains no hace falta bloquear a nadie. Aunque bloquees sigue pasando el mismo problema que aunque si bloquees el nodo y veas que existia, en el medio entre que devuelve la funcion alguien pudo haber borrado el nodo. 

# Ejercicio 4 
El metodo `edgeExists` hace el verify del metodo optimista, la idea es primero antes de hacer add o remove recorrer y bloquear el anterior y el siguiente. El problema es que antes de agarrar el anterior alguien pudo haberle eliminado la referencia, para asegurarnos de que nadie lo elimine antes, hay que volver a verificar de principio a fin que nadie haya eliminado las referencias. 