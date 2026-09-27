# Ejercicio 3 
Lo importante es que todos sigan el orden. 
Mientras todos vayan para un lado o para el otro todo bien. 
O vas de izquierda a derecha agarrando predecesor y actual 
O vas de derecha a izquierda tomando posterior y predecesor. 
Si un método sigue un orden y otro método sigue otro ahí vas a tener deadlock. 

El orden es tomar primero la cabeza, el siguiente, liberar la cabeza y tomar el siguiente del siguiente y liberar el siguiente. 

En el método contains no hace falta bloquear a nadie. Aunque bloquees sigue pasando el mismo problema que aunque si bloquees el nodo y veas que existía, en el medio entre que devuelve la función alguien pudo haber borrado el nodo. Así que tu punto de linealización es luego de determinar si está o no. 

# Ejercicio 4 
El método `edgeExists` hace el verify del método optimista, la idea es primero antes de hacer add o remove recorrer y bloquear el anterior y el siguiente. El problema es que antes de agarrar el anterior alguien pudo haberle eliminado la referencia, para asegurarnos de que nadie lo elimine antes, hay que volver a verificar de principio a fin que nadie haya eliminado las referencias. 

Al borrar un nodo `curr`, hay que conservar `curr.next` apuntando a su sucesor. Otro hilo puede haber leído `curr` antes del borrado y necesitar ese enlace para continuar recorriendo la lista. El borrado cambia `pred.next` para saltear `curr`, pero no modifica `curr.next`.

Si después de tomar los locks nos damos cuanta que `edgeExists` es falso, vamos a volver a recorrer la lista sin tomar locks hasta encontrar nuestros objetivos y bloquearlos y verificar. Potencialmente podrían todo el tiempo modificar los elementos que queremos modificar por lo que nuestro método no es wait-free, porque depende de las ejecuciones de los demás para poder terminar. 

# Ejercicio 5
El campo `marked` de la clase `Node` indica que el nodo ya fue removido lógicamente, aunque todavía no se haya actualizado el puntero del anterior para sacarlo de la lista.

Por más que podamos leer y escribir `marked` y `next` de forma atómica, para modificar la lista tenemos que hacer varias operaciones: verificar que la relación entre dos nodos siga siendo válida y después actualizar el enlace. Si no tomamos locks, otro thread podría cambiar esa relación entre la verificación y la actualización, y terminaríamos modificando una parte de la lista que ya cambió.

Para verificar que la relación `previo -> actual` siga existiendo, hay que bloquear los dos nodos y comprobar que `previo.next == actual` y que ninguno de los dos esté marcado. Así sabemos que la relación no va a cambiar antes de modificarla.

# Ejercicio 6
El método `compareAndSet` verifica que la referencia y la marca tengan los valores esperados y, si coinciden, cambia las dos cosas de forma atómica y devuelve `true`. Si alguna no coincide, no cambia nada y devuelve `false`. Esto puede pasar porque otro thread modificó el enlace o la marca desde que los leímos, así que hay que volver a buscar la posición y reintentar.

Un `volatile T` permite que los threads vean la referencia actualizada, pero no alcanza para comprobar su valor y cambiarlo en una sola operación ni para actualizar junto con ella una marca. `AtomicMarkableReference<T>` guarda la referencia y el booleano como un par, y permite leerlos y cambiarlos de forma atómica con `compareAndSet`.

El método `find` recorre la lista desde `head` buscando el primer nodo cuya clave sea mayor o igual a `value`. En cada paso lee el sucesor y la marca de `curr.next`. Si la marca está en `true`, significa que `curr` fue borrado lógicamente, no que su sucesor esté borrado. Entonces intenta saltear `curr` cambiando `pred.next` para que apunte a `succ`. Si ese cambio falla porque otro thread modificó la lista, vuelve a empezar desde `head`. Cuando encuentra un `curr` sin marcar con la clave buscada o una mayor, devuelve `pred` y `curr` en una `Window`. La marca en `next` permite distinguir un nodo borrado lógicamente de uno que todavía pertenece al conjunto, aunque aún no lo hayan sacado físicamente de la lista.
