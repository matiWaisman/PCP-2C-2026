# Ejercicio 1 

```
global semaphore puedeF = Semaphore(0)
global semaphore puedeC = Semaphore(0)

thread t1 {
    print(A)
    puedeF.release()
    print(B)
    puedeC.acquire()
    print(C)
}

thread t2 {
    print(E)
    puedeF.acquire()
    print(F)
    puedeC.release()
    print(G)
}
```

# Ejercicio 2 
Si le asignamos un timestamp a cada letra de cuando se tiene que mostrar como `Ts(l)`, se tiene que cumplir que: 

- `Ts(A) < Ts(C) < Ts(E) < Ts(R) < Ts(O)` o 
- `Ts(A) < Ts(C) < Ts(R) < Ts(E) < Ts(O)`

Para lograr que se cumpla al menos una, podemos hacer que `R` dependa de un semaforo que se activa luego de mostrar `C` y que `O` dependa de un semaforo que se activa luego de mostrar `E`. 

Notar que no hace falta que `E` dependa del mismo semaforo que `R` porque como asumimos que `print` es atomico, siempre se va a cumplir que `Ts(C) < Ts(E)`. Lo mismo ocurre entre `R` y `O`. 

```
global semaphore puedeC = Semaphore(0)
global semaphore puedeR = Semaphore(0)
global semaphore puedeO = Semaphore(0)

thread T1{
    puedeC.acquire()
    print(C)
    puedeR.release()
    print(E)
    puedeO.release()
}

thread T2{
    print(A)
    puedeC.release()
    puedeR.acquire()
    print(R)
    puedeO.acquire()
    print(O)
}
```

Como los `print` son atómicos, una vez que T1 ejecuta `puedeR.release()`, quedan dos instrucciones habilitadas para correr: `print(E)` en T1 y `print(R)` en T2 (esta última, apenas T2 termine su `acquire(puedeR)`, que en ese momento ya no bloquea). 

Si el scheduler le da el turno a T1 antes de que T2 llegue a su `print`, gana `E`; si le da el turno a T2 primero, gana `R`. En cualquiera de los dos casos, como el `print(R)` de T2 está antes que el `acquire(puedeO)` en el código de T2, `R` siempre termina impreso antes que `O`.

# Ejercicio 3 
```
global semaphore puedeI = Semaphore(0)
global semaphore puedeO = Semaphore(0)
global semaphore turnstile = Semaphore(0)

thread T1{
    print(R)
    puedeI.release()
    turnstile.acquire()
    print(OK)
}

thread T2{
    puedeI.acquire()
    print(I)
    puedeO.release()
    turnstile.acquire()
    print(OK)
}

thread T3{
    puedeO.acquire()
    print(O)
    turnstile.release(2)
    print(OK)
}
```

En este caso el patron que utilizamos es de barrera (un poco simplificado), todos se quedan esperando antes de dar el primer `OK` y los permisos del `OK` se liberan todos de una. No hace falta que `T3` tambien espere a la barrera porque sabemos que `print(o)` tiene que ser el ultimo print antes del primer ok. 

# Ejercicio 4 
```
global semaphore puedeF = semaphore(0)
global semaphore puedeH = semaphore(0)
global semaphore puedeC = semaphore(0)

thread T1{
    while (true){
        print(A)
        puedeF.release()
        print(B)
        puedeC.acquire()
        print(C)
        print(D)
    }
}

thread T2{
    while (true){
        print(E)
        puedeH.release()
        puedeF.acquire()
        print(F)
        print(G)
        puedeC.release()
    }
}

thread T3{
    while (true){
        puedeH.acquire()
        print(H)
        print(I)
    }
}
```

# Ejercicio 5 
## Punto A
```
global semaphore sA = Semaphore(1)
global semaphore sB = Semaphore(1)
thread T1{
    while(true){
        sA.acquire()
        print(A)
        sB.release()
    }
}

thread T2{
    while(true){
        sB.acquire()
        print(B)
        sA.release()
    }
}
```

## Punto B 
```
global semaphore sA = Semaphore(1)
global semaphore sB = Semaphore(0)
thread T1{
    while(true){
        sA.acquire()
        print(A)
        sB.release()
    }
}

thread T2{
    while(true){
        sB.acquire()
        print(B)
        sA.release()
    }
}
```
Este seria un caso de uso de semaforos split.

## Punto C
```
global semaphore sA = Semaphore(2)
global semaphore sB = Semaphore(0)
thread T1{
    while(true){
        sA.acquire(2)
        print(A)
        sB.release(2)
    }
}

thread T2{
    while(true){
        sB.acquire()
        print(B)
        sA.release()
    }
}
```

# Ejercicio 6 
El codigo está en [ej6.java](./ej6.java)

# Ejercicio 7
## Punto A 
Asumo que hay un solo lugar donde se dejan todos los discos (un rack central) y que la cantidad total de discos es acotada.

Los recursos compartidos son: las cuatro máquinas (cada una independiente de las demás) y el rack de discos, que es un único recurso compartido por todos los clientes y todas las máquinas. Los agentes activos son los clientes del gimnasio, cada uno con una rutina propia: una lista finita de pasos, donde cada paso indica qué máquina usar y cuántos discos necesita para ese ejercicio. La rutina puede repetir máquinas o no usar alguna.

Para sincronizar el acceso propongo un mutex por máquina (cuatro en total, uno por aparato) y un semáforo contador para el rack de discos, inicializado con la cantidad total de discos disponibles. Cada vez que un cliente necesita discos, hace un `acquire` de tantos permisos como discos precise; al devolverlos, hace el `release` correspondiente.

Para cada paso de su rutina, un cliente:
1. Espera a que la máquina que necesita esté libre (`acquire` del mutex de esa máquina).
2. Una vez que la tiene, pide la cantidad de discos que su rutina indica (`acquire` del semáforo de discos).
3. Realiza el ejercicio.
4. Descarga la máquina, devolviendo los discos al rack (`release` del semáforo de discos).
5. Libera la máquina (`release` del mutex), quedando disponible para el siguiente cliente.

El orden importa en ambas puntas. Al principio, primero se ocupa la máquina y recién después se piden los discos, para que nadie quede cargado de discos mientras espera turno por una máquina. Al final, primero se devuelven los discos y después se libera la máquina, respetando la norma del gimnasio de que un aparato no debe quedar "usable" por otro cliente hasta que esté descargado.

Este orden fijo de adquisición —siempre máquina antes que discos— es lo que evita el deadlock: como todos los clientes piden los recursos en el mismo orden, no puede darse una espera circular donde un cliente tenga la máquina esperando discos mientras otro tiene discos esperando esa misma máquina.

Con lo que si hay que tener cuidado es con la obtencion de los discos, si es atomica haciendo un `acquires(cantDiscosNecesarias)` no va a haber problema, pero si se hace con un ciclo de `acquires` puede haber un deadlock. Un caso puede ser que dos personas quieren usar maquinas distintas que ambas estan vacias, si cada uno necesita 3 discos y hay 3 disponibles, puede ocurrir que cada uno agarre 2 y sea interrumpido y el otro agarre un disco y se queden en deadlock. 


## Punto B
El codigo está en [ej7.java](./ej7.java)

## Punto C 
En caso de que todo el tiempo llegue gente a una maquina que un usuario este esperando, y "tenga mala suerte" una persona no podría usar nunca una maquina. Para solucionar esto podriamos hacer que el semaforo sea fuerte y simule una fila, algo que es razonable en un gimnasio. 

Otro problema que puede haber puede ser a la hora de agarrar discos, si hay pocos discos y las 4 personas se estan peleando por usar los pocos discos que hay, tambien podria hacer que una persona que ya es su turno en la maquina puede llegar a tener "mala suerte" y nunca pueda agarrar discos para la maquina, porque siempre les roba los discos otro, para solucionar esto tambien podemos hacer que el semaforo sea fuerte, y nuevamente tenemos otra fila para agarrar discos. 

# Ejercicio 8 
TODO 

# Ejercicio 9 
## Punto A 

```
global semaphore s0 = Semaphore(1)
global semaphore s1 = Semaphore(0)
global semaphore[2] s = [s0, s1]
global int constaActual = 0
global int capacidadTransbordador = N

thread transbordador(){

}

thread persona(int costa){
    s[costa].acquire()

}
```

# Ejercicio 11 
## Punto A 
```
global int personasEnElBaño = 0

global semaphore mutexPB = Semaphore(1)
global semaphore baños = Semaphore(8)
global semaphore acceso = Semaphore(1)


thread persona(){
    mutexPB.acquire()
    personasEnElBaño += 1
    if(personasEnElBaño == 1){
        acceso.acquire()
    }
    mutexPB.release()
    baños.acquire()
    print("Haciendo pichin")
    baños.release()
    mutexPB.acquire()
    personasEnElBaño -= 1 
    if(personasEnElBaño == 0){
        acceso.release()
    }
    mutexPB.release()
}

thread limpieza(){
    acceso.acquire()
    print("Limpiando")
    acceso.release()
}
```

## Punto B 
Si el prioridad de limpieza tiene prioridad, lo que podemos hacer es agregar un molinete FIFO, que hace que nadie pueda entrar al baño si esta el personal de limpieza esperando, por lo que cuando llega el personal de limpieza va a esperar a que todos los que esten adentro terminen y despues va a entrar el. 

```
global int personasEnElBaño = 0

global semaphore molinete = Semaphore(1, True)
global semaphore mutexPB = Semaphore(1)
global semaphore baños = Semaphore(8)
global semaphore acceso = Semaphore(1)


thread persona(){
    molinete.acquire()
    molinete.release()
    mutexPB.acquire()
    personasEnElBaño += 1
    if(personasEnElBaño == 1){
        acceso.acquire()
    }
    mutexPB.release()
    baños.acquire()
    print("Haciendo pichin")
    baños.release()
    mutexPB.acquire()
    personasEnElBaño -= 1 
    if(personasEnElBaño == 0){
        acceso.release()
    }
    mutexPB.release()
}

thread limpieza(){
    molinete.acquire()
    acceso.acquire()
    print("Limpiando")
    molinete.release()
    acceso.release()
}
```

## Punto C
```
global int personasEnElBaño = 0

global semaphore mutexPB = Semaphore(1)
global semaphore baños = Semaphore(8)
global semaphore acceso = Semaphore(1)
global semaphore molinete = Semaphore(1, True)


thread persona(){
    mutexPB.acquire()
    personasEnElBaño += 1
    if(personasEnElBaño == 1){
        acceso.acquire()
    }
    mutexPB.release()
    baños.acquire()
    print("Haciendo pichin")
    baños.release()
    mutexPB.acquire()
    personasEnElBaño -= 1 
    if(personasEnElBaño == 0){
        acceso.release()
    }
    mutexPB.release()
}

thread limpieza(){
    acceso.acquire()
    print("Limpiando")
    acceso.release()
}
```
