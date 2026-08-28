# Ejercicio 1
## Punto A

Viendo primero la semántica de los hilos por separado: 

- Si ignoramos la declaración e inicialización de las variables globales:

    - $\llbracket T_1 \rrbracket = \lambda\sigma.\sigma[y \mapsto \sigma(x)]$

    - $\llbracket T_2 \rrbracket = \lambda\sigma.\sigma[x \mapsto \sigma(y)]$

- Si aplicamos cada hilo por separado al estado inicial
$\sigma_0=[x\mapsto1,y\mapsto2]$:

    - $\llbracket T_1 \rrbracket\sigma_0 = \sigma_0[y\mapsto1] = [x\mapsto1,y\mapsto1]$

    - $\llbracket T_2 \rrbracket\sigma_0 = \sigma_0[x\mapsto2] = [x\mapsto2,y\mapsto2]$

Ahora, viendo la ejecución paralela de ambos hilos:

- $\llbracket T_1\parallel T_2\rrbracket=\lambda\sigma.\sigma[y\mapsto\sigma(x)]\oplus\lambda\sigma.\sigma[x\mapsto\sigma(y)]\oplus\lambda\sigma.\sigma[x\mapsto\sigma(y),y\mapsto\sigma(x)]$

Para el diagrama, cada thread va a tener dos instrucciones, en la primera lee la variable que va a asignarle a la otra y luego la escribe. Para eso estan `tX`, que representa la variable en la que T1 guarda el valor de `x`, y `tY`, que representa la variable en la que T2 guarda el valor de `y`. 

![Diagrama de transición de estados de T1 paralelo T2](./images/p1e1a.png)

## Punto B 

Viendo primero la semántica de los hilos por separado: 
- $\llbracket T_1 \rrbracket = \lambda\sigma.\sigma[y \mapsto \sigma(x) + 1]$
- $\llbracket T_2 \rrbracket = \lambda\sigma.\sigma[x \mapsto \sigma(y) + 1]$

Para el caso de la ejecución paralela: 
- $\llbracket T_1\parallel T_2\rrbracket=\lambda\sigma.\sigma[x \mapsto \sigma(y) + 1, y \mapsto \sigma(x) + 1] \oplus \lambda\sigma.\sigma[x \mapsto \sigma(y) + 1, y \mapsto \sigma(y) + 2] \oplus \lambda\sigma.\sigma[x \mapsto \sigma(x) + 2, y \mapsto \sigma(x) + 1]$

![Diagrama de transición de estados de T1 paralelo T2](./images/p1e1b.png)

## Punto C 

Semántica por separado: 
- $\llbracket T_1\rrbracket=\lambda\sigma.\operatorname{if}\ \llbracket x<1\rrbracket\sigma\ \operatorname{then}\bot\ \operatorname{else}\sigma$
- $\llbracket T_2 \rrbracket = \lambda\sigma.\sigma[x \mapsto 1]$

Ejecución Paralela: 
- $\llbracket T_1\parallel T_2\rrbracket=\displaystyle\bigoplus_{n\in\mathbb{N}_0}\lambda\sigma.\sigma[x\mapsto1,y\mapsto\sigma(y)+n]\oplus\bot$

Para el diagrama de transición de estados asumo que todas las operaciones son atomicas para disminuir la complejidad del grafico, y tambien asumo que es weakly fair. 

![Diagrama de transición de estados de T1 paralelo T2](./images/p1e1c.png)

Deberían haber infinitos estados, pero la idea es la del grafico. 

# Ejercicio 2

Separo los programas en sus operaciones atómicas:

```yaml
global n = 0

thread T1
    local tmp
    p1: do K times
        p2: tmp = n
        p3: n = tmp + 1

thread T2
    local tmp
    q1: do K times
        q2: tmp = n
        q3: n = tmp + 1
```

## Punto A

Para que el valor final de `n` sea `2K`, se puede ejecutar primero un hilo durante sus `K` iteraciones y luego el otro. También pueden alternarse, siempre que cada incremento se complete antes de que el otro hilo lea `n`.

| T1             | T2             | Estado                       |
|----------------|----------------|------------------------------|
| `do K times`   |                | `p:p2`                       |
| `tmp = n`      |                | `p.tmp = 0; p:p3`            |
| `n = tmp + 1`  |                | `n = 1; p:p1`                |
| ...            |                | ...                          |
| `do K times`   |                | `p:p2`                       |
| `tmp = n`      |                | `p.tmp = K-1; p:p3`          |
| `n = tmp + 1`  |                | `n = K; p:p1`                |
| `do K times`   |                | `p:−`                        |
|                | `do K times`   | `q:q2`                       |
|                | `tmp = n`      | `q.tmp = K; q:q3`            |
|                | `n = tmp + 1`  | `n = K+1; q:q1`              |
|                | ...            | ...                          |
|                | `do K times`   | `q:q2`                       |
|                | `tmp = n`      | `q.tmp = 2K-1; q:q3`         |
|                | `n = tmp + 1`  | `n = 2K; q:q1`               |
|                | `do K times`   | `q:−`                        |

## Punto B
Para que el valor final sea `k`, lo que tendría que pasar es que se vayan ejecutando sincronizadamente T1 y T2 dentro del ciclo, donde siempre leen ambos el valor de tmp antes de actualizar la variable, y escribe uno y despues el otro. 

| T1            | T2            | Estado              |
|---------------|---------------|---------------------|
| `do K times`  |               | `p:p2`              |
|               | `do K times`  | `q:q2`              |
| `tmp = n`     |               | `p.tmp = 0; p:p3`   |
|               | `tmp = n`     | `q.tmp = 0; q:q3`   |
| `n = tmp + 1` |               | `n = 1; p:p1`       |
|               | `n = tmp + 1` | `n = 1; q:q1`       |
| ...           | ...           | ...                 |
| `do K times`  |               | `p:p2`              |
|               | `do K times`  | `q:q2`              |
| `tmp = n`     |               | `p.tmp = k-1; p:p3` |
|               | `tmp = n`     | `q.tmp = k-1; q:q3` |
| `n = tmp + 1` |               | `n = k; p:p1`       |
|               | `n = tmp + 1` | `n = k; q:q1`       |
| `do K times`  |               | `p:-`               |
|               | `do K times`  | `q:-`               |

## Punto C 
Si, puede ocurrir para el siguiente entramado con `K=3`

| T1            | T2            | Estado            |
|---------------|---------------|-------------------|
| `do K times`  |               | `p:p2`            |
| `tmp = n`     |               | `p.tmp = 0; p:p3` |
|               | `do K times`  | `q:q2`            |
|               | `tmp = n`     | `q.tmp = 0; q:q3` |
|               | `n = tmp + 1` | `n=1; q:q1`       |
|               | `do K times`  | `q:q2`            |
|               | `tmp = n`     | `q.tmp = 1; q:q3` |
|               | `n = tmp + 1` | `n=2; q:q1`       |
| `n = tmp + 1` |               | `n = 1; p:p1`     |
|               | `do K times`  | `q:q2`            |
|               | `tmp = n`     | `q.tmp = 1; q:q3` |
| `do K times`  |               | `p:p2`            |
| `tmp = n`     |               | `p.tmp = 1; p:p3` |
| `n = tmp + 1` |               | `n = 2; p:p1`     |
| `do K times`  |               | `p:p2`            |
| `tmp = n`     |               | `p.tmp = 2; p:p3` |
| `n = tmp + 1` |               | `n = 3; p:p1`     |
| `do K times`  |               | `p:-`             |
|               | `n = tmp + 1` | `n=2; q:q1`       |
|               | `do K times`  | `q:-`             |

# Ejercicio 3 
Siendo que el programa es: 

```yaml
global x = 1
global y = 2
global z = 3

thread T1
    local tX
    p1: tX = x
    p2: y = tX

thread T2
    local tY
    q1: tY = y
    q2: z = tY

thread T3
    local tZ
    r1: tZ = z
    r2: x = tZ
```

Los resultados posibles pueden ser: 
| x | y | z |
|---|---|---|
| 1 | 1 | 1 |
| 2 | 2 | 2 |
| 3 | 3 | 3 |
| 3 | 1 | 2 |
| 3 | 1 | 1 |
| 2 | 1 | 2 |
| 3 | 3 | 2 |

Donde los primeros 3 surguen de ejecuciones secuenciales en algun orden de cada thread. 
El cuarto es en el cual los 3 threads leen antes de que el primer escriba. 
Y los ultimos 3 donde primero se escribe una variable y luego las otras dos leen antes de que la primera escriba y luego escriben las dos. 

# Ejercicio 4 
## Punto A
Si tenemos `N` threads y `K` instrucciones, cada ejecucion del programa tiene que ejecutar exactamente `N * K` lineas para terminar la ejecucion del programa. 

Como todas las ejecuciones posibles toman exactamente `N * K` pasos y hay que sumarle uno en el diagrama de estados por el estado inicial, entonces la cota minima de estados en el diagrama es de `N * K + 1`.

Esta cota es exacta si solo existe un posible orden de ejecucion de los distintos threads. En cambio si los threads son independientes entre si en ese caso hay una mayor cantidad de estados porque en cada nodo podemos llegar a tener `N` posibles aristas a nodos distintos. 

Una cota superior va a ser (K + 1)^N. Cada estado del diagrama queda 
determinado por el vector de program counters de los N threads, donde cada 
pc_i puede tomar independientemente uno de K + 1 valores posibles (de 0 a K, 
sin terminar y terminado inclusive). Como los threads son independientes, 
todas las combinaciones de estos N valores son alcanzables, dando 
(K + 1)^N estados posibles en total.


## Punto B
Como cada camino tiene longitud `NK + 1`

TODO

# Ejercicio 5
## Punto A 
El programa no es correcto. Contraejemplo: sea f una funcion con una unica 
raiz entera r. Considerar el siguiente entrelazado:

1. T1 ejecuta y encuentra la raiz, sale del while y termina, dejando found = true.
2. Recien ahora T2 ejecuta sus primeras lineas: local j = 1; found = false, 
   pisando el true que dejo T1.
3. T2 entra al while(!found) y empieza a buscar en los negativos indefinidamente.

Como f tiene una unica raiz positiva, T2 no termina, violando la definicion de correctitud dada.

## Punto B 
El programa no es correcto. Contraejemplo: sea f una funcion con una unica 
raiz entera r. Considerar el siguiente entrelazado: 

1. T2 ejecuta hasta la comparacion no atomica `while (!found)`, por ahora lo unico que hizo es leer el valor viejo de found. 
2. T1 ejecuta y encuentra la raiz, sale del while y termina, dejando found = true.
3. T2 continua ya que en su lectura `found` era falso, por lo que luego lo sobreescribe con `found = (f(j) == 0)`, por lo que va a quedar buscando en los negativos indefinidamente. 

Como f tiene una unica raiz positiva, T2 no termina, violando la definicion de correctitud dada.

## Punto C
Si el scheduler es fair y deja ejecutar a ambos procesos, entonces el programa es correcto. 

Ya que si al menos uno de los dos threads encuentra la raiz, no puede ocurrir un lost update en el que el otro pise el valor con falso para hacer que uno o ambos siga buscando. Asi que como mucho el proceso que no encuentra la raiz va a tener que ejecutar una iteracion mas. Pero ambos van a terminar si existe una raiz. 

El unico problema que puede haber es si el scheduler no es fair y hay una sola raiz y no deja ejecutar al proceso que explora del lado de la raiz. 

# Ejercicio 6
## Punto A 
Un 2 en la salida puede aparecer 0 o 1 veces: 
- Si se ejecuta todo `T2` antes que `T1` entonces no se va a entrar al while, por lo que va a aparecer 0 veces.
- Si se termina de ejecutar `T2` antes de leer el valor de `n` para hacer el print, entonces se va a printear un 2 y luego se va a salir del while. 

## Punto B 
El 1 puede aparecer entre cero y una cantidad infinita de veces: 
- Va a aparecer cero veces si se ejecuta todo `T2` antes que `T1`, por lo que no se va a entrar al while.
- Puede aparecer entre 1 e infinitas veces si se ejecuta la primer linea de `T2` y luego se queda ejecutando el ciclo de `T1`. Si el scheduler no fuera fair podria quedarse colgado el thread `T1` printeando 1. 

## Punto C
Puede no mostrarse nada por pantalla si se ejecuta `T2` en su totalidad antes que se ejecute `T1` por primera vez. 

# Ejercicio 7
## Punto A 
Si, existe un interleaving en el que el loop `T1` se ejecute exactamente una vez, esto puede ocurrir si se ejecuta primero la totalidad de `T1` y luego se ejecute `T2`. 

## Punto B
Si separamos el codigo del programa en las instrucciones atomicas: 

- `p1`/`q1`: Lectura de n en una variable local `tmp`.
- `p2`/`q2`: Comparacion y desicion del loop contra `tmp`.
- `p3`/`q3`: Escritura nueva en `n`

Que los ciclos de `T1` y `T2` terminen dependen del orden en el que ejecute los hilos el scheduler, en un caso en el que los programas se ejecuten ciclicamente de la forma: 

| T1            | T2            | Estado            |
|---------------|---------------|-------------------|
| `p1`          |               | `p.tmp = 0; p:p2` |
| `p2`          |               | `p:p3`            |
|               | `q1`          | `q.tmp = 0; q:q2` |
|               | `q2`          | `q:q3`            |
| `n = tmp + 1` |               | `n=1; p:p1`       |
|               | `n = tmp - 1` | `n=0; q:q1`       |
| `p1`          |               | `p.tmp = 0; p:p2` |
| `p2`          |               | `p:p3`            |
|               | `q1`          | `q.tmp = 0; q:q2` |
|               | `q2`          | `q:q3`            |
| `n = tmp + 1` |               | `n=1; p:p1`       |
|               | `n = tmp - 1` | `n=0; q:q1`       |

Esto haria que ni `T1` ni `T2` terminen. 

# Ejercicio 8
## Punto A 
Los posibles valores de n pueden ir desde 0 hasta el menos infinito. 

Los valores del programa pueden valer entre -1 y menos infinito, en caso de que se empiece ejecutando el segundo thread y si la condicion de `n ==0` requiere de dos pasos, el primero en el que se lee el valor de la variable compartida `n` y se guarda en una variable local, y otro en el que se efectua el if, en el medio el scheduler puede poner a correr al otro proceso y mientras tanto disminuye el valor global de `n`. Otra posibilidad sea que el scheduler corte al segundo thread justo antes de poner la flag en falso o durante algun paso no atomico. 

## Punto B 
La ejecucion del programa puede no terminar en caso de que se actualize el valor de `n` antes de que se llegue a la primera lectura de `n` dentro del segundo thread. En ese caso el programa nunca va a terminar. 

# Ejercicio 9 
Para ver si un algoritmo resuelve el problema de la exclusion mutua lo que hay que comprobar es que cumpla con las siguientes tres propiedades: 

- Mutex: en todo momento, a lo sumo un thread está en la sección crítica.
- Ausencia de deadlock/livelock: no puede darse una situación donde el sistema quede trabado y ningún thread que quiere entrar a la sección crítica pueda progresar (ni bloqueado esperando indefinidamente, ni dando vueltas sin avanzar).
- Garantía de entrada: todo thread que quiere entrar a la sección crítica, eventualmente entra (no puede quedar postergado para siempre, incluso si otros threads sí logran entrar repetidamente).

Este algoritmo no cumple con la primer propiedad, ambos threads pueden estar en la seccion critica simultanteamente. Como las lineas `np = nq + 1` y `nq = np + 1` no son atomicas y minimo requieren de dos pasos, uno en el que se lee el valor de la variable compartida y otro en el que se escribe el valor nuevo en la variable compartida, puede ocurrir que ambos lleguen a la linea del `while` con el mismo numero de prioridad, haciendo que ninguno se quede loopeando y logrando que ambos entren a la seccion critica a la vez y ninguno se quede loopeando. Ya que no se van a cumplir las condiciones para tener que esperar porque la prioridad del otro va a ser distinta a cero y tampoco se cumple que uno tiene mayor prioridad que el otro, ya que tienen la misma. 

Si separamos el codigo en las siguientes operaciones atomicas: 

```
global np = 0
global nq = 0

thread p                          thread q
  while(true){                      while(true){
p1:   tmp_p = nq                q1:     tmp_q = np
p2:   np = tmp_p + 1            q2:     nq = tmp_q + 1
p3:   while (nq != 0            q3:     while (np != 0
          && np > nq) {}                    && nq > np) {}
p4:   // seccion critica        q4:     // seccion critica
p5:   np = 0                    q5:     nq = 0
  }                                  }
```

Llegamos a un interleaving en el que ambos threads estan en la zona critica a la vez:

| p                                    | q                                    | Estado              |
|--------------------------------------|--------------------------------------|---------------------|
| `tmp_p = nq`                         |                                      | `p.tmp_p = 0; p:p2` |
|                                      | `tmp_q = np`                         | `q.tmp_q = 0; q:q2` |
| `np = tmp_p + 1`                     |                                      | `np = 1; p:p3`      |
|                                      | `nq = tmp_q + 1`                     | `nq = 1; q:q3`      |
| `while (nq != 0 && np > nq)` (false) |                                      | `p:p4`              |
|                                      | `while (np != 0 && nq > np)` (false) | `q:q4`              |
| `p4`                                 |                                      | `p:p5`              |
|                                      | `q4`                                 | `q:q5`              |

# Ejericicio 10 
## Punto A 

En este algoritmo no se cumple la propiedad de Mutex, si varios procesos primero leen `local turno = turnos` `turno` va a tener valor cero y luego es reemplazado por el scheduler por otro proceso este tambien va a tener el turno cero, luego ambos van a entrar a la seccion critica a la vez. 

Si nuestro codigo es: 
```
global actual = 0
global turnos = 0

PedirTurno(){
n1:   local turno = turnos
n2:   turnos = turnos + 1
n3:   return turno
}

LiberarTurno(){
n4:   actual = actual + 1
n5:   turnos = turnos - 1
}

//SECCION NO CRITICA
n6:   local miturno = PedirTurno()
n7:   while (actual != miturno) {}
n8:   //SECCION CRITICA
n9:   LiberarTurno()
      //SECCION NO CRITICA
```


| p                              | q                              | r                              | Estado                            |
|--------------------------------|--------------------------------|--------------------------------|-----------------------------------|
| `local miturno = PedirTurno()` |                                |                                | `p:p1`                            |
| `local turno = turnos`         |                                |                                | `p.turno = 0; p:p2`               |
|                                | `local miturno = PedirTurno()` |                                | `q:q1`                            |
|                                | `local turno = turnos`         |                                | `q.turno = 0; q:q2`               |
|                                |                                | `local miturno = PedirTurno()` | `r:r1`                            |
|                                |                                | `local turno = turnos`         | `r.turno = 0; r:r2`               |
| `turnos = turnos + 1`          |                                |                                | `turnos = 1; p.miturno = 0; p:p3` |
|                                | `turnos = turnos + 1`          |                                | `turnos = 2; q.miturno = 0; q:q3` |
|                                |                                | `turnos = turnos + 1`          | `turnos = 3; r.miturno = 0; r:r3` |
| `while (actual != miturno)`    |                                |                                | `p:p8`                            |
|                                | `while (actual != miturno)`    |                                | `q:q8`                            |
|                                |                                | `while (actual != miturno)`    | `r:r8`                            |
| `p8`                           |                                |                                | `p:p9`                            |
|                                | `q8`                           |                                | `q:q9`                            |
|                                |                                | `r8`                           | `r:r9`                            |

## Punto B 
Si `PedirTurno` y `LiberarTurno` fueran atomicas, nunca podria ocurrir que varios procesos tengan el mismo turno. Por lo que la propiedad de mutex se cumple. Y la propiedad de ausencia de Deadlock se va a cumplir porque siempre para algun proceso se va a cumplir que `actual` va a ser igual a su turno. Tambien se va a cumplir la garantia de entrada asumiendo que en la seccion critica ningun thread se puede bloquear, haciendo que no avance el turno. 

# Ejercicio 11
En el algoritmo propuesto no se cumple mutex. Si tenemos por ejemplo 4 threads y tanto el primero como el tercero encienden su flag y van al `while` a ver si pueden entrar o no a la seccion critica, como ni el thread 2 ni el 4 encendieron sus flags, la condicion `flag [otro]` va a ser falsa, por lo que los dos pueden entrar al ciclo a la vez. 

Siendo el codigo: 

```
global flag[n] = {false, false, ..., false}
global turno = 0

thread(id) {
      //SECCION NO CRITICA
n1:   flag[id] = true
n2:   local otro = (id + 1) % n
n3:   turno = otro
n4:   while (flag[otro] && turno == otro) {}
n5:   //SECCION CRITICA
n6:   flag[id] = false
      //SECCION NO CRITICA
}
```

| p                                        | q | r                                        | s | Estado                                |
|------------------------------------------|---|------------------------------------------|---|---------------------------------------|
| `flag[id] = true`                        |   |                                          |   | `flag=[true,false,false,false]; p:p2` |
| `local otro = (id + 1) % n`              |   |                                          |   | `p.otro = 1; p: p3`                   |
| `turno = otro`                           |   |                                          |   | `turno = 1; p: p4`                    |
|                                          |   | `flag[id] = true`                        |   | `flag=[true,false,true,false]; r:r2`  |
|                                          |   | `local otro = (id + 1) % n`              |   | `r.otro = 3; r: r3`                   |
|                                          |   | `turno = otro`                           |   | `turno = 3; r: r4`                    |
| `flag[id] = true`                        |   |                                          |   | `flag=[true,false,false,false]; p:p2` |
| `while (flag[otro] && turno == otro) {}` |   |                                          |   | `p:p5`                                |
| `p5`                                     |   |                                          |   | `p:p6`                                |
|                                          |   | `while (flag[otro] && turno == otro) {}` |   | `r:r5`                                |
|                                          |   | `r5`                                     |   | `r:r6`                                |

No pueden ocurrir ni deadlock ni imposibilidad de entrada a un thread en particular. 

Como mucho puede quedarse esperando un unico thread a la vez, el que coincida su id con el anterior a el que fue el ultimo en sobreescribir `turno`. Pero una vez que termine (asumiendo que siempre termina la seccion critica) tambien va a poder entrar. 

# Ejercicio 12 
## Punto A
El algoritmo no resuelve el problema de exlusion mutua porque puede ocurrir un Deadlock. 

Si al menos dos procesos setean su flag en True antes de que el primero entre a la seccion critica, ambos van a quedarse trabados en el while, luego todos los demas que quieran entrar tambien se van a quedar stuckeados en el while. Como ocurre deadlock tambien puede ocurrir procesos no puedan entrar nunca. 

Lo que si cumple el algoritmo es mutex. 

Supongamos, por absurdo, que `p` y `q` entran a la sección crítica "a la vez". Para que eso pase, la llamada a `algunVerdadero(p)` que hace `p` tiene que devolver `false` (en particular, en algún instante $t_{p2}$ lee `flag[q]` y la ve en `false`), y análogamente `algunVerdadero(q)` tiene que devolver `false` (en algún instante $t_{q2}$ lee `flag[p]` y la ve en `false`).

Pero sabemos:
- $t_{p1}$ (cuando `p` hace `flag[p] = true`) es anterior a $t_{p2}$ (la lectura de `flag[q]` durante el scan de `p`), porque `p` primero setea su propia flag y recién después arranca el while.
- Análogamente, $t_{q1} < t_{q2}$.

Para que `p` no vea a `q`, hace falta $t_{p2} < t_{q1}$ (p lee antes de que q escriba). Para que `q` no vea a `p`, hace falta $t_{q2} < t_{p1}$.

Encadenando todo: $t_{p2} < t_{q1} < t_{q2} < t_{p1} < t_{p2}$ — un ciclo, absurdo. Entonces es imposible que ambos scans devuelvan `false` simultáneamente respecto del otro. Por lo tanto, no pueden estar los dos en la sección crítica a la vez.

## Punto B 
Si `algunVerdadero` fuera atomico, aun asi se puede producir el deadlock si mas de un proceso setea su flag en true antes que el siguiente entre a la seccion critica. 

# Ejercicio 13 
TODO 
Dice que hay que eliminar un `j < id` y no hay ninguno en el codigo. Y difiere mucho vs el bakery de la teorica. 

# Ejercicio 14 
El codigo no cumple con la ausencia de deadlocks, por lo que no cumple exclusion mutua. El problema esta en el uso de la variable compartida `turno`. En caso de que a un thread se le asigne un ticket mayor a cero, nunca va a poder pasar del while `while ( turno != miturno ) {}` porque turno siempre vale cero. Esto puede ocurrir facilmente si un primer thread ejecuta el primer `fetch-and-add`, en ese caso su turno va a ser el 0 y el ticket va a valer 1. Luego, todo proceso que llame a `fetch-and-add` antes de que el primer thread salga de la seccion critica se le va a asignar un ticket distinto a cero, por lo que va a quedar bloqueado. 

Para arreglarlo, lo que habria que hacer es reemplazar la linea antes de la seccion critica por `fetch-and-add(turno, miturno, 1)`. De esta manera, cuando un proceso sale de la seccion critica, mueve el turno hacia arriba permitiendo ejecutar al siguiente thread. Como el `fetch-and-add` es atomico, no puede ocurrir que dos threads tengan el mismo turno. 

TODO Demo mas rigurosa de porque cumple todo 

# Ejercicio 15
## Punto A
Si la operacion `tomarFlag` no es atomica, la propuesta anterior no resuelve el problema de la exclusión mutua. La solución propuesta no cumple con no tener deadlock. 

Como la línea `flag[mia] = !flag[otro]` no es atomica y requiere por lo menos dos instrucciones, uno de leer el valor del array y otro para actualizarlo. Esto puede hacer que ambos seteen su flag como verdadero si el otro tambien tiene en falso su flag. 

Si descomponemos el codigo en: 
```
global flag[0..1]={false, false}


tomarFlag (mia , otro ) {
n1:    local otro = flag[otro]
n2:    flag [mia] = !otro
}

thread T0 {
n3: while (!flag [0]){
n4:     tomarFlag (0 ,1)
    }
n5: // SECCION CRITICA
n6: flag [0]= false
}
```

Podemos armar una traza como: 

| p                         | q                         | Estado                        |
|---------------------------|---------------------------|-------------------------------|
| `while (!flag[0])`        |                           | `flag = [false, false]; p:p4` |
| `tomarFlag(0,1)`          |                           | `p:p2`                        |
|                           | `while (!flag[1])`        | `flag = [false, false]; q:q4` |
|                           | `tomarFlag(1,0)`          | `q:q2`                        |
| `local otro = flag[otro]` |                           | `p.otro = 0; p:p3`            |
|                           | `local otro = flag[otro]` | `q.otro = 0; q:q3`            |
| `flag[mia] = !otro`       |                           | `flag=[true,false]; p:p5`     |
|                           | `flag[mia] = !otro`       | `flag=[true,true]; q:q5`      |
| `while (!flag[0])`        |                           | `flag = [true, true]; p:p5`   |
|                           | `while (!flag[1])`        | `flag = [true, true]; q:q5`   |
| `p5`                      |                           | `p:p6`                        |
|                           | `q5`                      | `q:q6`                        |

Entonces en este interleaving los dos threads entran en la seccion critica a la vez. 

## Punto B 
Si `tomarFlag` fuera atomica, no puede ocurrir un interleaving como el anterior, porque siempre exactamente un thread va a ganar la carrera y va a setear su bandera en true, y como va a ser el primero el anterior va a quedar en false. 





