# Ejercicio 1 
## Punto A 

[Solución en Haskell](ej1.hs)

## Punto B 

Sí, tanto los clientes como los agentes pueden sufrir starvation. STM en Haskell
ejecuta las transacciones de forma optimista, registrando sus lecturas y escrituras
en un log. Antes de confirmar sus cambios, valida que la ejecución sea consistente.
Si otra transacción confirmó una modificación que invalida lo leído, los cambios
tentativos se descartan y la transacción debe ejecutarse nuevamente.

No se garantiza que cada transacción consiga confirmar sus cambios: un thread
puede perder indefinidamente frente a otros que sí progresan. Reintentar muchas
veces y finalmente completar no constituye starvation; debe existir una ejecución
en la que el thread nunca logre completar la operación.

Una traza posible para un cliente es la siguiente:

1. El cliente `C` lee su cola y comprueba que hay espacio.
2. Antes de que `C` confirme, otro cliente encola su llamada y confirma la
   modificación de esa misma cola.
3. La transacción de `C` resulta inválida y debe reiniciarse.
4. Un agente toma una llamada y libera espacio. Al reintentar `C`, otro cliente
   vuelve a confirmar un encolado entre la lectura de `C` y su validación.
5. Con un flujo continuo de clientes y atenciones, este patrón puede repetirse
   indefinidamente: otros clientes logran encolarse, pero `C` nunca lo consigue.

Una traza posible para un agente es la siguiente:

1. Hay una llamada urgente en espera y ninguna común. El agente `A` lee las colas
   e intenta retirar la urgente.
2. Antes de que `A` confirme, otro agente `B` retira esa llamada y confirma su
   transacción, invalidando la de `A`.
3. `A` debe reiniciar. Llega una nueva urgente y `B` termina su atención anterior.
4. En el siguiente intento, `A` vuelve a leer la cola, pero `B` vuelve a retirar la
   llamada antes de que `A` confirme.
5. Si este patrón se repite indefinidamente, `B` continúa atendiendo llamadas
   mientras `A` nunca logra tomar una, aunque siga ejecutando intentos.

La implementación no incorpora un mecanismo que garantice equidad entre las
transacciones y evite estos abortos indefinidos. La regla de las tres urgentes
regula qué tipo de llamada se selecciona, pero no garantiza que cada cliente o
agente consiga confirmar sus transacciones. Estas trazas muestran starvation de
un thread, no un deadlock: otros threads siguen avanzando.

## Punto C
Asumiendo que ambos centros trabajan a la misma velocidad, cada cliente en su codigo antes de encolarse podria ver la capacidad de las dos colas de su prioridad en ambos centros, y anotarse en el centro que tiene menor cantidad de clientes encolados. 