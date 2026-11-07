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