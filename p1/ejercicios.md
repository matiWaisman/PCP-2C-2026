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