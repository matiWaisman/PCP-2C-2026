# Ejercicio 1 

- En granularidad gruesa: 
    - Si ocurre que un elemento `b` que no pertenece al conjunto tiene el mismo hash que un elemento `a` que si pertenece, cuando en el código del `add` se hace `if (key == curr.key) return false`, la comparación ahora debería ser `if (key == curr.key && curr.value == o.value) return false`. 
    - En el `remove` no habría que cambiar nada porque la desición de eliminar es con un if por el valor del elemento. 
- 