# Ejercicio 1 
```mermaid
flowchart TD
    S1["(p1, q1, x=0)"]
    S2["(p2 [loc=0], q1, x=0)"]
    S3["(p1, q2 [loc=0], x=0)"]
    S4["(fin, q1, x=1)"]
    S5["(p1, fin, x=1)"]
    S6["(p2 [loc=0], q2 [loc=0], x=0)"]
    S7["(fin, q2 [loc=1], x=1)"]
    S8["(p2 [loc=1], fin, x=1)"]
    S9["(fin, q2 [loc=0], x=1)"]
    S10["(p2 [loc=0], fin, x=1)"]
    S11["(fin, fin, x=2)"]
    S12["(fin, fin, x=1)"]

    S1 -->|p1| S2
    S1 -->|q1| S3
    S2 -->|p2| S4
    S2 -->|q1| S6
    S3 -->|q2| S5
    S3 -->|p1| S6
    S4 -->|q1| S7
    S5 -->|p1| S8
    S6 -->|p2| S9
    S6 -->|q2| S10
    S7 -->|q2| S11
    S8 -->|p2| S11
    S9 -->|q2| S12
    S10 -->|p2| S12

    style S1 stroke:#3b82f6,stroke-width:3px
    style S6 stroke:#ef4444,stroke-width:3px
    style S11 stroke:#22c55e,stroke-width:3px
    style S12 stroke:#ef4444,stroke-width:3px
```
# Ejercicio 2
El LTS de `p` y `q` por separado: 

```mermaid
flowchart LR
    subgraph P
        direction LR
        P1["p1"] -->|acquire_a.p| P2["p2"]
        P2 -->|acquire_b.p| P4["p4"]
        P4 -->|release_b.p| P5["p5"]
        P5 -->|release_a.p| PF["fin"]
    end

    subgraph Q
        direction LR
        Q1["q1"] -->|acquire_b.q| Q2["q2"]
        Q2 -->|acquire_a.q| Q4["q4"]
        Q4 -->|release_a.q| Q5["q5"]
        Q5 -->|release_b.q| QF["fin"]
    end

    style P1 stroke:#3b82f6,stroke-width:3px
    style Q1 stroke:#3b82f6,stroke-width:3px
```

Y el LTS de la composición es: 
```mermaid
flowchart TD
    S0["(p1, q1)"]
    S1["(p2, q1)"]
    S2["(p1, q2)"]
    S3["(p4, q1)"]
    S4["(p2, q2)"]
    S5["(p1, q4)"]
    S6["(p5, q1)"]
    S7["(p1, q5)"]
    S8["(fin, q1)"]
    S9["(p5, q2)"]
    S10["(p2, q5)"]
    S11["(p1, fin)"]
    S12["(fin, q2)"]
    S13["(p2, fin)"]
    S14["(fin, q4)"]
    S15["(p4, fin)"]
    S16["(fin, q5)"]
    S17["(p5, fin)"]
    S18["(fin, fin)"]

    S0 -->|acquire_a.p| S1
    S0 -->|acquire_b.q| S2
    S1 -->|acquire_b.p| S3
    S1 -->|acquire_b.q| S4
    S2 -->|acquire_a.p| S4
    S2 -->|acquire_a.q| S5
    S3 -->|release_b.p| S6
    S5 -->|release_a.q| S7
    S6 -->|release_a.p| S8
    S6 -->|acquire_b.q| S9
    S7 -->|acquire_a.p| S10
    S7 -->|release_b.q| S11
    S8 -->|acquire_b.q| S12
    S9 -->|release_a.p| S12
    S10 -->|release_b.q| S13
    S11 -->|acquire_a.p| S13
    S12 -->|acquire_a.q| S14
    S13 -->|acquire_b.p| S15
    S14 -->|release_a.q| S16
    S15 -->|release_b.p| S17
    S16 -->|release_b.q| S18
    S17 -->|release_a.p| S18

    style S0 stroke:#3b82f6,stroke-width:3px
    style S4 stroke:#ef4444,stroke-width:4px
    style S18 stroke:#22c55e,stroke-width:3px
```

Los caminos más cortos hacia un deadlock tienen largo 2 y son las ejecuciones en las que p toma `a` (`acquire_a.p`) y q toma `b` (`acquire_b.q`), en cualquiera de los dos órdenes. Ambas llegan al estado `(p2, q2)`, donde p tiene `a` y espera `b`, y q tiene `b` y espera `a`, por lo que ninguna acción está habilitada.

1. $(p_1,q_1) \xrightarrow{\text{acquire\_a.p}} (p_2,q_1) \xrightarrow{\text{acquire\_b.q}} (p_2,q_2)$
2. $(p_1,q_1) \xrightarrow{\text{acquire\_b.q}} (p_1,q_2) \xrightarrow{\text{acquire\_a.p}} (p_2,q_2)$

Suponiendo que el thread `q()` pasa a usar `tryLock()`, su código va a ser: 

```c
void q() {
q1: b.lock()
q2: if(a.tryLock()){
q3:     usar()
q4:     a.unlock()
q5: }
q6: b.unlock()
}
```

El LTS de `q` va a ser: 

```mermaid
flowchart LR
    Q1["q1"] -->|acquire_b.q| Q2["q2"]
    Q2 -->|trylock_a_true.q| Q4["q4"]
    Q2 -->|trylock_a_false.q| Q5["q5"]
    Q4 -->|release_a.q| Q5
    Q5 -->|release_b.q| QF["fin"]

    style Q1 stroke:#3b82f6,stroke-width:3px
```

El LTS de la composición va a ser: 

```mermaid
flowchart TD
    S0["(p1, q1)"]
    S1["(p2, q1)"]
    S2["(p1, q2)"]
    S3["(p4, q1)"]
    S4["(p2, q2)"]
    S5["(p1, q4)"]
    S6["(p5, q1)"]
    S7["(p2, q5)"]
    S8["(p1, q5)"]
    S9["(fin, q1)"]
    S10["(p5, q2)"]
    S11["(p2, fin)"]
    S12["(p1, fin)"]
    S13["(fin, q2)"]
    S14["(p5, q5)"]
    S15["(p4, fin)"]
    S16["(fin, q4)"]
    S17["(fin, q5)"]
    S18["(p5, fin)"]
    S19["(fin, fin)"]

    S0 -->|acquire_a.p| S1
    S0 -->|acquire_b.q| S2
    S1 -->|acquire_b.p| S3
    S1 -->|acquire_b.q| S4
    S2 -->|acquire_a.p| S4
    S2 -->|trylock_a_true.q| S5
    S3 -->|release_b.p| S6
    S4 -->|trylock_a_false.q| S7
    S5 -->|release_a.q| S8
    S6 -->|release_a.p| S9
    S6 -->|acquire_b.q| S10
    S7 -->|release_b.q| S11
    S8 -->|acquire_a.p| S7
    S8 -->|release_b.q| S12
    S9 -->|acquire_b.q| S13
    S10 -->|release_a.p| S13
    S10 -->|trylock_a_false.q| S14
    S11 -->|acquire_b.p| S15
    S12 -->|acquire_a.p| S11
    S13 -->|trylock_a_true.q| S16
    S14 -->|release_a.p| S17
    S14 -->|release_b.q| S18
    S15 -->|release_b.p| S18
    S16 -->|release_a.q| S17
    S17 -->|release_b.q| S19
    S18 -->|release_a.p| S19

    style S0 stroke:#3b82f6,stroke-width:3px
    style S4 stroke:#f59e0b,stroke-width:3px
    style S19 stroke:#22c55e,stroke-width:3px
```

No hay deadlock porque todos los estados tienen al menos una transición saliente. El único que no tiene ninguna es `fin` (que se puede pensar que tiene un loop en sí mismo). Tampoco se pueden producir ciclos dentro del LTS.

Si el orden de adquisición es igual entre los dos threads, el LTS va a ser: 

```mermaid
flowchart TD
    A0["(p1, q1)"]
    A1["(p2, q1)"]
    A2["(p4, q1)"]
    A3["(p5, q1)"]
    A4["(fin, q1)"]
    A5["(fin, q2)"]
    A6["(fin, q4)"]
    A7["(fin, q5)"]
    B1["(p1, q2)"]
    B2["(p1, q4)"]
    B3["(p1, q5)"]
    B4["(p1, fin)"]
    B5["(p2, fin)"]
    B6["(p4, fin)"]
    B7["(p5, fin)"]
    F["(fin, fin)"]

    A0 -->|acquire_a.p| A1
    A1 -->|acquire_b.p| A2
    A2 -->|release_b.p| A3
    A3 -->|release_a.p| A4
    A4 -->|acquire_a.q| A5
    A5 -->|acquire_b.q| A6
    A6 -->|release_b.q| A7
    A7 -->|release_a.q| F

    A0 -->|acquire_a.q| B1
    B1 -->|acquire_b.q| B2
    B2 -->|release_b.q| B3
    B3 -->|release_a.q| B4
    B4 -->|acquire_a.p| B5
    B5 -->|acquire_b.p| B6
    B6 -->|release_b.p| B7
    B7 -->|release_a.p| F

    style A0 stroke:#3b82f6,stroke-width:3px
    style F stroke:#22c55e,stroke-width:3px
```

No hay deadlock porque todo estado distinto de (`fin`, `fin`) tiene exactamente una transición saliente y no se produce ningún ciclo. Entonces de todo nodo se puede llegar al final.

# Ejercicio 3 
Para el caso sin buffer para cada thread el LTS va a ser: 

```mermaid
flowchart LR
    subgraph MAIN
        direction LR
        M1["m1"] -->|C1| M2["m2"]
        M2 -->|C2| MF["fin"]
    end

    subgraph ANON
        direction LR
        A1["a1"] -->|C2| A2["a2"]
        A2 -->|C1| AF["fin"]
    end

    style M1 stroke:#3b82f6,stroke-width:3px
    style A1 stroke:#3b82f6,stroke-width:3px
```

Como los canales no tienen buffer, enviar y recibir son un único evento: el envío por un canal queda bloqueado hasta que otra goroutine esté lista para recibir por ese mismo canal, y lo mismo al revés. Por eso `C1` y `C2` son acciones compartidas entre MAIN y ANON, y solo ocurren si las dos las ofrecen a la vez.

Por lo tanto al hacer la composicion el LTS solo se compone de un estado porque las transiciones no se pueden hacer al necesitar que los dos a la vez hagan la accion de enviar y recibir sobre el mismo buffer. Por lo tanto es un deadlock ya que nos quedamos trabados en un estado que no es el final. 

```mermaid
flowchart TD
    S0["(m1, a1)<br/>DEADLOCK"]

    style S0 stroke:#ef4444,stroke-width:4px
```

Ahora si el Buffer tiene capacidad uno, no es necesario que el emisor y receptor esten sincronizados para hacer el envio, por lo tanto los LTS de cada componente por separado va a ser: 

**MAIN**

```mermaid
flowchart LR
    M1["m1"] -->|Envia C1| M2["m2"]
    M2 -->|Recibe C2| MF["fin"]

    style M1 stroke:#3b82f6,stroke-width:3px
```

**ANON**

```mermaid
flowchart LR
    A1["a1"] -->|Envia C2| A2["a2"]
    A2 -->|Recibe C1| AF["fin"]

    style A1 stroke:#3b82f6,stroke-width:3px
```

**C1** (canal `c1` con buffer de 1)

```mermaid
flowchart LR
    C1V["vacío"] -->|Envia C1| C1L["lleno"]
    C1L -->|Recibe C1| C1V

    style C1V stroke:#3b82f6,stroke-width:3px
```

**C2** (canal `c2` con buffer de 1)

```mermaid
flowchart LR
    C2V["vacío"] -->|Envia C2| C2L["lleno"]
    C2L -->|Recibe C2| C2V

    style C2V stroke:#3b82f6,stroke-width:3px
```

La composicion va a ser: 

```mermaid
flowchart TD
    S0["(m1, a1, vacío, vacío)"]
    S1["(m2, a1, lleno, vacío)"]
    S2["(m1, a2, vacío, lleno)"]
    S3["(m2, a2, lleno, lleno)"]
    S4["(fin, a2, lleno, vacío)"]
    S5["(m2, fin, vacío, lleno)"]
    S6["(fin, fin, vacío, vacío)"]

    S0 -->|Envia C1| S1
    S0 -->|Envia C2| S2
    S1 -->|Envia C2| S3
    S2 -->|Envia C1| S3
    S3 -->|Recibe C2| S4
    S3 -->|Recibe C1| S5
    S4 -->|Recibe C1| S6
    S5 -->|Recibe C2| S6

    style S0 stroke:#3b82f6,stroke-width:3px
    style S6 stroke:#22c55e,stroke-width:3px
```

Esta versión no puede quedar bloqueada. Si un proceso intenta recibir antes de que el otro haya enviado, queda esperando (no tiene transición habilitada), pero el otro proceso siempre puede avanzar: cada goroutine hace su envío antes de su recepción, y como los buffers arrancan vacíos, el envío nunca se bloquea. Cuando el dato queda en el buffer, la recepción pendiente se habilita. En el LTS de la composición, el único estado sin transiciones salientes es `(fin, fin, vacío, vacío)`, que es la terminación correcta, así que no existe ningún deadlock.

# Ejercicio 4 
# Ejercicio 4

### 1. $\square\Diamond q$

*Desde cualquier posición, más adelante (o en ella misma) vale $q$. Es decir, $q$ ocurre infinitas veces.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **No vale.**
    - $q$ vale por última vez en $i = 3$. Desde $i \geq 4$ solo vale $p$, así que $\Diamond q$ falla ahí y, por lo tanto, $\square\Diamond q$ falla en $i = 1$.
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **Vale.**
    - $q$ vale en todas las posiciones pares. En las impares vale en la posición siguiente.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **Vale.**
    - En $i = 1$ vale $q$. En $i = 2$ vale en la posición siguiente. Desde $i = 3$ vale siempre.
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **Vale.**
    - $q$ vale en $i = 3, 6, 9, \dots$, es decir, infinitas veces. Desde cualquier posición hay un $q$ más adelante.

### 2. $\Diamond\square p$

*Eventualmente siempre se cumple $p$.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **Vale.**
    - A partir de $i = 4$ siempre vale $p$.
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **No vale.**
    - En las posiciones pares no vale $p$, por lo tanto no existe un $i$ tal que a partir de ese $i$ siempre valga $p$.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **Vale.**
    - A partir de $i = 3$ siempre vale $p$.
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **No vale.**
    - No existe $i$ tal que a partir de ese $i$ valga $p$ para todas las posiciones, ya que hay posiciones en las que no se cumple nada.

### 3. $p \,\mathbf{U}\, q$

*$q$ vale en algún momento y, hasta entonces, vale $p$.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **Vale.**
    - En $i = 1$ vale $p$ y $q$ aparece en $i = 2$. En $i = 2$ y $i = 3$ vale $q$ ahí mismo. Entonces $p \,\mathbf{U}\, q$ vale para $1 \leq i \leq 3$. Desde $i = 4$ ya no vale, porque $q$ no vuelve a ocurrir.
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **Vale.**
    - En las posiciones pares vale $q$ ahí mismo. En las impares vale $p$ y en la posición siguiente vale $q$. Entonces vale para todo $i$.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **Vale.**
    - En $i = 1$ vale $q$ ahí mismo, así que no se le pide nada a $p$. En $i = 2$ falla (no vale ni $p$ ni $q$), pero eso no afecta a $i = 1$. Desde $i = 3$ vale $q$ en todas las posiciones.
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **No vale.**
    - En $i = 1$ vale $p$, pero en $i = 2$ no vale ni $p$ ni $q$, y $q$ todavía no apareció. La fórmula falla en $i = 1$. En general, solo vale en $i = 3, 6, 9, \dots$ (donde vale $q$ ahí mismo).

### 4. $\square(p \rightarrow \Diamond q)$

*Cada vez que vale $p$, en ese momento o más adelante vale $q$.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **No vale.**
    - Una vez que se entra al ciclo de $p$ nunca mas ocurre $q$.
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **Vale.**
    - Siempre luego de una $p$ hay una $q$.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **Vale.**
    - Siempre que se cumple $p$ tambien se cumple $q$.
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **Vale.**
    - Siempre luego de una $p$ en dos pasos se cumple $q$.

### 5. $\bigcirc\bigcirc q$

*En la posición $i = 3$ vale $q$.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **Vale.**
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **No vale.**
    - En $i = 3$ vale $p$ y no $q$.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **Vale.**
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **Vale.**

### 6. $\square(p \,\mathbf{U}\, q)$

*Siempre vale que $q$ ocurre en algun momento, y hasta entonces se cumple $p$.*

- $\sigma_1 = \{p\}\,\{p,q\}\,\{q\}\,(\{p\})^\omega$: **No vale.**
    - A partir de $i = 4$ no se vuelve a cumplir $q$.
- $\sigma_2 = (\{p\}\,\{q\})^\omega$: **Vale.**
    - Para toda posicion o se cumple $q$ o se cumple $p$ y luego se va a cumplir $q$.
- $\sigma_3 = \{q\}\,\emptyset\,(\{p,q\})^\omega$: **No vale.**
    - No vale para $i = 2$, por lo que no vale, en esa posicion no vale ni $p$ ni $q$.
- $\sigma_4 = \{p\}\,\emptyset\,(\{p,q\}\,\{p\}\,\emptyset)^\omega$: **No vale.**
    - No vale para $i = 2$ porque no se cumple ni $p$ ni $q$ en esa posicion, por lo que no se cumple.