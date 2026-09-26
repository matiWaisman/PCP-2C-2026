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

Para lograr que se cumpla al menos una, podemos hacer que `R` dependa de un semáforo que se activa luego de mostrar `C` y que `O` dependa de un semáforo que se activa luego de mostrar `E`. 

Notar que no hace falta que `E` dependa del mismo semáforo que `R` porque como asumimos que `print` es atómico, siempre se va a cumplir que `Ts(C) < Ts(E)`. Lo mismo ocurre entre `R` y `O`. 

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

En este caso el patrón que utilizamos es de barrera (un poco simplificado), todos se quedan esperando antes de dar el primer `OK` y los permisos del `OK` se liberan todos de una. No hace falta que `T3` también espere a la barrera porque sabemos que `print(o)` tiene que ser el último print antes del primer ok. 

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
Este sería un caso de uso de semáforos split.

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
El código está en [ej6.java](./ej6.java)

# Ejercicio 7
## Punto A 
Asumo que hay un solo lugar donde se dejan todos los discos (un rack central) y que la cantidad total de discos es acotada.

Los recursos compartidos son: las cuatro máquinas (cada una independiente de las demás) y el rack de discos, que es un único recurso compartido por todos los clientes y todas las máquinas. Los agentes activos son los clientes del gimnasio, cada uno con una rutina propia: una lista finita de pasos, donde cada paso indica qué máquina usar y cuántos discos necesita para ese ejercicio. La rutina puede repetir máquinas o no usar alguna.

Para sincronizar el acceso propongo un mutex por máquina (cuatro en total, uno por aparato) y un semáforo contador para el rack de discos, inicializado con la cantidad total de discos disponibles. Cada vez que un cliente necesita discos, hace un `acquire` de tantos permisos como discos precise; al devolverlos, hace el `release` correspondiente. Es importante notar que esta solución funciona porque hacer `semaphore.acquire(n)` en java es atomico, o se toman los n permisos o ninguno, si no fuera así podría ocurrir un deadlock. 

Para cada paso de su rutina, un cliente:
1. Espera a que la máquina que necesita esté libre (`acquire` del mutex de esa máquina).
2. Una vez que la tiene, pide la cantidad de discos que su rutina indica (`acquire` del semáforo de discos).
3. Realiza el ejercicio.
4. Descarga la máquina, devolviendo los discos al rack (`release` del semáforo de discos).
5. Libera la máquina (`release` del mutex), quedando disponible para el siguiente cliente.

El orden importa en ambas puntas. Al principio, primero se ocupa la máquina y recién después se piden los discos, para que nadie quede cargado de discos mientras espera turno por una máquina. Al final, primero se devuelven los discos y después se libera la máquina, respetando la norma del gimnasio de que un aparato no debe quedar "usable" por otro cliente hasta que esté descargado.

Este orden fijo de adquisición —siempre máquina antes que discos— es lo que evita el deadlock: como todos los clientes piden los recursos en el mismo orden, no puede darse una espera circular donde un cliente tenga la máquina esperando discos mientras otro tiene discos esperando esa misma máquina.

Con lo que sí hay que tener cuidado es con la obtención de los discos, si es atómica haciendo un `acquires(cantDiscosNecesarias)` no va a haber problema, pero si se hace con un ciclo de `acquires` puede haber un deadlock. Un caso puede ser que dos personas quieren usar máquinas distintas que ambas están vacías, si cada uno necesita 3 discos y hay 3 disponibles, puede ocurrir que cada uno agarre 2 y sea interrumpido y el otro agarre un disco y se queden en deadlock. 


## Punto B
El código está en [ej7.java](./ej7.java)

Asumo que `rackDiscos.acquire(cantidadDiscos);` se hace de manera atómica: O se hacen los acquire de `cantidadDiscos` o no se hace nada, no equivale a un ciclo de `acquire(1)`, porque si fuera así podría ocurrir un deadlock. 

Si el `acquires(n)` equivale a hacer `n` iteraciones de `acquire(1)` para salvar un caso de deadlock como el que hablé en el inciso anterior habría que encerrar la acción de `acquire` del rack dentro de un mutex, para que solo pueda haber una persona intentando agarrar discos a la vez. Como todas las personas van a terminar eventualmente de usar la máquina tarde o temprano se van a liberar los discos suficientes para que cada persona pueda agarrar y no ocurra un deadlock. 

## Punto C 
En caso de que todo el tiempo llegue gente a una máquina que un usuario esté esperando, y "tenga mala suerte" una persona no podría usar nunca una máquina. Para solucionar esto podríamos hacer que el semáforo sea fuerte y simule una fila, algo que es razonable en un gimnasio. 

Otro problema que puede haber puede ser a la hora de agarrar discos, si hay pocos discos y las 4 personas se están peleando por usar los pocos discos que hay, también podría hacer que una persona que ya es su turno en la máquina puede llegar a tener "mala suerte" y nunca pueda agarrar discos para la máquina, porque siempre les roba los discos otro, para solucionar esto también podemos hacer que el semáforo sea fuerte, y nuevamente tenemos otra fila para agarrar discos. 

# Ejercicio 8 
Este problema suena mucho a productores y consumidores con un buffer acotado, donde los consumidores tienen que consumir de a dos a la vez, y para los productores la cota para producir es dinámica respecto de la cantidad de productores que hay. 

Por lo que dice el enunciado, siempre puedo asumir que por lo menos hay 2 o más productores. 

```java 
global int cantidadBolitas = 0;

global semaphore notEmpty = Semaphore(0);
global semaphore slots = Semaphore(0);
global semaphore mutexCantidadElementos = Semaphore(1);

thread generador(){
    void empezarAGenerar(){
        slots.release();
    }

    void dejarDeGenerar(){
        slots.acquire(); // Le resto un permiso que nunca más va a volver
    }

    int generar(){
        slots.acquire();
        mutexCantidadElementos.acquire();
        cantidadBolitas += 1;
        if(cantidadBolitas % 2 == 0){
            notEmpty.release();
        }
        mutexCantidadElementos.release();
        return 1;
    }
}

thread consumidor(){
    int consumir(){
        notEmpty.acquire();
        mutexCantidadElementos.acquire();
        cantidadBolitas -= 2; 
        slots.release(2);
        mutexCantidadElementos.release();
        return 1;
    }
}
```


# Ejercicio 9 
## Punto A 

```java
global semaphore s0 = Semaphore(0);
global semaphore s1 = Semaphore(0);
global semaphore[2] s = [s0, s1];
global semaphore mutexPersonasEnTransbordador = Semaphore(1);
global semaphore puedePartir = Semaphore(0);
global semaphore puedeTerminar = Semaphore(0);
global semaphore puedenBajar = Semaphore(0);
global int costaActual = 0;
global int capacidadTransbordador;
global int cantidadPersonasEnTransbordador = 0;



thread transbordador(int N){
    void transbordador(){
        mutexPersonasEnTransbordador.acquire();
        capacidadTransbordador = N; 
        s[0].release(capacidadTransbordador);
        mutexPersonasEnTransbordador.release();
    }

    void iniciarViaje(){
        puedePartir.acquire();
    }

    void finalizarViaje(){
        puedenBajar.release(capacidadTransbordador);
        puedeTerminar.acquire();
        costaActual = (costaActual + 1) % 2;
        s[costaActual].release(capacidadTransbordador);
    }
}

thread persona(int costa){
    void subirBote(int costaOrigen){
        s[costaOrigen].acquire();
        mutexPersonasEnTransbordador.acquire();
        cantidadPersonasEnTransbordador += 1;
        if(cantidadPersonasEnTransbordador == capacidadTransbordador){
            puedePartir.release();
        }
        mutexPersonasEnTransbordador.release();
    }

    void bajarseBote(int costaDestino){
        puedenBajar.acquire();
        mutexPersonasEnTransbordador.acquire();
        cantidadPersonasEnTransbordador -= 1; 
        if(cantidadPersonasEnTransbordador == 0){
            puedeTerminar.release();
        }
        mutexPersonasEnTransbordador.release();
    }
}
```
## Punto B 
Asumo que por más que ahora la gente puede subir y bajar concurrentemente, aun así para empezar el nuevo viaje tienen que bajarse los `N` del viaje anterior y subir los `N` del viaje nuevo. 

```java
global semaphore s0 = Semaphore(0);
global semaphore s1 = Semaphore(0);
global semaphore[2] s = [s0, s1];
global semaphore mutexPersonasSubieron = Semaphore(1);
global semaphore puedePartir = Semaphore(0);
global semaphore puedenBajar = Semaphore(0);
global int costaActual = 0;
global int capacidadTransbordador;
global int cantidadPersonasSubieron = 0; 

thread transbordador(int N){
    void transbordador(){
        capacidadTransbordador = N; 
        s[0].release(capacidadTransbordador);
    }

    void iniciarViaje(){
        puedePartir.acquire();
    }

    void finalizarViaje(){
        costaActual = (costaActual + 1) % 2;
        mutexPersonasSubieron.acquire();
        cantidadPersonasSubieron = 0;
        mutexPersonasSubieron.release();
        puedenBajar.release(capacidadTransbordador);
    }
}

thread persona(int costa){
    void subirBote(int costaOrigen){
        s[costaOrigen].acquire();
        mutexPersonasSubieron.acquire();
        cantidadPersonasSubieron += 1;
        if(cantidadPersonasSubieron == capacidadTransbordador){
            puedePartir.release();
        }
        mutexPersonasSubieron.release();
    }

    void bajarseBote(int costaDestino){
        puedenBajar.acquire();
        s[costaActual].release(); // Libero un lugar
    }
}
```

# Ejercicio 10 
Dentro de los posibles lugares de carga y descarga, las máquinas van a tener los ids de 0 a 7, la plataforma de recepción el id de 8, y la de entrega id 9. 

En cualquier momento puede ir un vehículo a la plataforma de descarga o de carga a hacer descarga o carga y no importa si ya hay otros haciendo lo mismo. 

Asumo que tengo acceso a un objeto dummy `FakeSemaphore()` que implementa `acquire` y `release` pero que ambas son no bloqueantes y realmente no hacen nada. De esta manera no hace falta poner ifs que si el id de la carga o descarga es uno en particular haga algo distinto y siempre se haga lo mismo independientemente de en que estación de carga o descarga se va a operar.

Asumo también que cada auto tiene un id del 0 al 3 indicando cuál es el cual puede consultar haciendo `currentThread.id()` y lo mismo con las máquinas/ estaciones de carga y descarga que tienen ids del 0 al 9. 

```java 
global List<Semaphore> autoPuedeDescargarEnMaquina = new List<Semaphore>();
// mutex: controla que un auto pueda ocupar el slot de descarga de la máquina (auto -> máquina)

global List<Semaphore> maquinaRecibioMateriaPrima = new List<Semaphore>();
// señal: el auto avisa a la máquina que ya dejó materia prima para procesar (auto -> máquina)

global List<Semaphore> autoPuedeCargarDeMaquina = new List<Semaphore>();
// mutex: controla que un auto pueda ocupar el slot de carga de la máquina (auto -> máquina)

global List<Semaphore> maquinaTerminoDeProcesar = new List<Semaphore>();
// señal: la máquina avisa a los autos que ya hay producto refinado listo para retirar (máquina -> auto)

// Posiciones 0 a 7: las 8 máquinas procesadoras
for(int i = 0; i < 8; i++){
    autoPuedeDescargarEnMaquina.push(Semaphore(1));  // arranca libre
    maquinaRecibioMateriaPrima.push(Semaphore(0));   // arranca sin nada para procesar
    autoPuedeCargarDeMaquina.push(Semaphore(1));     // arranca libre
    maquinaTerminoDeProcesar.push(Semaphore(0));     // arranca sin producto listo
}

// Posiciones 8 y 9: plataforma de recepcion (8) y de entrega (9) - todo fake
autoPuedeDescargarEnMaquina.push(FakeSemaphore());
autoPuedeDescargarEnMaquina.push(FakeSemaphore());
maquinaRecibioMateriaPrima.push(FakeSemaphore());
maquinaRecibioMateriaPrima.push(FakeSemaphore());
autoPuedeCargarDeMaquina.push(FakeSemaphore());
autoPuedeCargarDeMaquina.push(FakeSemaphore());
maquinaTerminoDeProcesar.push(FakeSemaphore());
maquinaTerminoDeProcesar.push(FakeSemaphore());

thread maquina(int identificador){
    int idMaquina = identificador;
    
    void descargar(){
        maquinaRecibioMateriaPrima[currentThread.id()].acquire();
        
    }

    void procesar(){

        // Hace algo
    }

    void cargar(){
        maquinaTerminoDeProcesar[currentThread.id()].release();
    }
}

thread vehiculo(int identificador){
    int idVehiculo = identificador;

    void cargar(int idDestino){
        autoPuedeCargarDeMaquina[idDestino].acquire();
        maquinaTerminoDeProcesar[idDestino].acquire();
        autoPuedeCargarDeMaquina[idDestino].release();
    }

    void descargar(int idDestino){
        autoPuedeDescargarEnMaquina[idDestino].acquire();
        maquinaRecibioMateriaPrima[idDestino].release();
        autoPuedeDescargarEnMaquina[idDestino].release();
    }
}
```
# Ejercicio 11
El ejercicio es similar al de lectores escritores, donde la gente que usa el baño son los lectores, y el personal de limpieza el escritor.  
## Punto A 
```java
global int personasEnElBaño = 0;

global semaphore mutexPB = Semaphore(1);
global semaphore baños = Semaphore(8);
global semaphore acceso = Semaphore(1);

thread persona(){
    mutexPB.acquire();
    personasEnElBaño += 1;
    if(personasEnElBaño == 1){
        acceso.acquire();
    }
    mutexPB.release();
    baños.acquire();
    print("Haciendo pichin");
    baños.release();
    mutexPB.acquire();
    personasEnElBaño -= 1; 
    if(personasEnElBaño == 0){
        acceso.release();
    }
    mutexPB.release();
}

thread limpieza(){
    acceso.acquire();
    print("Limpiando");
    acceso.release();
}
```

## Punto B 
Si el prioridad de limpieza tiene prioridad, lo que podemos hacer es agregar un molinete FIFO, que hace que nadie pueda entrar al baño si está el personal de limpieza esperando, por lo que cuando llega el personal de limpieza va a esperar a que todos los que estén adentro terminen y después va a entrar él. 

```java
global int personasEnElBaño = 0;

global semaphore molinete = Semaphore(1, True);
global semaphore mutexPB = Semaphore(1);
global semaphore baños = Semaphore(8);
global semaphore acceso = Semaphore(1);


thread persona(){
    molinete.acquire();
    molinete.release();
    mutexPB.acquire();
    personasEnElBaño += 1;
    if(personasEnElBaño == 1){
        acceso.acquire();
    }
    mutexPB.release();
    baños.acquire();
    print("Haciendo pichin");
    baños.release();
    mutexPB.acquire();
    personasEnElBaño -= 1; 
    if(personasEnElBaño == 0){
        acceso.release();
    }
    mutexPB.release();
}

thread limpieza(){
    molinete.acquire();
    acceso.acquire();
    print("Limpiando");
    molinete.release();
    acceso.release();
}
```

## Punto C
```java 
global int personasEnElBaño = 0;

global semaphore mutexPB = Semaphore(1);
global semaphore baños = Semaphore(8);
global semaphore acceso = Semaphore(1);
global semaphore molinete = Semaphore(1, true);

thread persona(){
    baños.acquire();
    molinete.acquire();
    molinete.release();
    mutexPB.acquire();
    personasEnElBaño += 1;
    if(personasEnElBaño == 1){
        acceso.acquire();
    }
    mutexPB.release();
    print("Haciendo pichin");
    mutexPB.acquire();
    personasEnElBaño -= 1; 
    if(personasEnElBaño == 0){
        acceso.release();
    }
    mutexPB.release();
    baños.release();
}

thread limpieza(){
    molinete.acquire();
    acceso.acquire();
    print("Limpiando");
    molinete.release();
    acceso.release();
}
```
# Ejercicio 12
## Punto A
Asumo que las dos posibles direcciones/ destinos están representadas con el cero y el uno. 
```java 
global int[] cantidadAutosCruzando = {0, 0};
global semaphore[] mutexAutosCruzando = {Semaphore(1), Semaphore(1)};
global semaphore acceso = Semaphore(1);

thread auto(){
    void entrarPuente(int direccionOrigen){
        mutexAutosCruzando[direccionOrigen].acquire();
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.acquire();
        }
        // Si estoy acá es o porque ya tenia el acceso de antes mi lado o lo acabo de ganar 
        cantidadAutosCruzando[direccionOrigen] += 1;
        mutexAutosCruzando[direccionOrigen].release();
    }

    void cruzarPuente(){
        print("Cruzando el puente");
    }

    void salirPuente(int direccionOrigen){
        mutexAutosCruzando[direccionOrigen].acquire();
        cantidadAutosCruzando[direccionOrigen] -= 1;
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.release();
        }
        mutexAutosCruzando[direccionOrigen].release();
    }
}
```

La idea sería que cada auto ejecute un dependiendo de la dirección en la que parte: 
```
entrarPuente(origen);
cruzarPuente();
salirPuente(origen);
```

El problema que tiene esta primer solución es que si todo el tiempo entran autos de una dirección, deja en starvation a los autos del otro lado. Para solucionar esto se puede agregar a `entrarPuente` un turnstile que lo que haga es hacer que quien cruza el puente lo determina el orden en el que llegó al principio. Es un trade off entre más serialización por menos inanición versus menos inanición y más serialización. 

```java 
global int[] cantidadAutosCruzando = {0, 0};
global semaphore[] mutexAutosCruzando = {Semaphore(1), Semaphore(1)};
global semaphore acceso = Semaphore(1);
global semaphore turnstile = Semaphore(1, true);

thread auto(){
    void entrarPuente(int direccionOrigen){
        turnstile.acquire();
        mutexAutosCruzando[direccionOrigen].acquire();
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.acquire();
        }
        // Si estoy acá es o porque ya tenia el acceso de antes mi lado o lo acabo de ganar 
        cantidadAutosCruzando[direccionOrigen] += 1;
        mutexAutosCruzando[direccionOrigen].release();
        turnstile.release();
    }

    void cruzarPuente(){
        print("Cruzando el puente");
    }

    void salirPuente(int direccionOrigen){
        mutexAutosCruzando[direccionOrigen].acquire();
        cantidadAutosCruzando[direccionOrigen] -= 1;
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.release();
        }
        mutexAutosCruzando[direccionOrigen].release();
    }
}
```

## Punto B 
```java 
global int[] cantidadAutosCruzando = {0, 0};
global semaphore[] mutexAutosCruzando = {Semaphore(1), Semaphore(1)};
global semaphore autosCruzando = Semaphore(3);
global semaphore acceso = Semaphore(1);

thread auto(){
    void entrarPuente(int direccionOrigen){
        mutexAutosCruzando[direccionOrigen].acquire();
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.acquire();
        }
        // Si estoy acá es o porque ya tenia el acceso de antes mi lado o lo acabo de ganar 
        cantidadAutosCruzando[direccionOrigen] += 1;
        mutexAutosCruzando[direccionOrigen].release();
        autosCruzando.acquire();
    }

    void cruzarPuente(){
        print("Cruzando el puente");
    }

    void salirPuente(int direccionOrigen){
        autosCruzando.release();
        mutexAutosCruzando[direccionOrigen].acquire();
        cantidadAutosCruzando[direccionOrigen] -= 1;
        if(cantidadAutosCruzando[direccionOrigen] == 0){
            acceso.release();
        }
        mutexAutosCruzando[direccionOrigen].release();
    }
}
```

## Punto C 
La solución propuesta no es libre de inanición, si empiezan cruzando los autos de una dirección, y constantemente entran autos de esa dirección nunca van a poder pasar los autos de la otra dirección. Para solucionar esto se podría agregar un turnstile como en el inciso A para que no haya inanición, a costo de por ejemplo un caso que si de un lado tenemos 3 autos pero llegó primero uno del otro lado, va a pasar primero ese auto en vez de los 3 autos del otro lado que son más.

El límite de que solo puedan cruzar 3 a la vez no impide que haya inanición por el modo en el que se maneja el cambio de lado.