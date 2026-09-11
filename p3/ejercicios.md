# Ejercicio 1 
## Punto A 
Si la disciplina utilizada es la de Hoare "signal y espera urgente", esta solución no cumple con el requerimiento de que siempre el orden sea `antes→importante→despues`.

Una traza que lo rompe es una en la que el thread que ejecuta `antesydespues` llega primero al monitor. En ese caso, entra, hace la parte de `antes`, y cuando llega a `permiso.signal()`, como todavía no hay nadie esperando en la condition `permiso`, el signal no tiene ningún efecto (no despierta a nadie porque no hay nadie en la cola). El thread simplemente sigue ejecutando dentro del monitor y hace `despues`, sin que `importante` se haya ejecutado nunca.

Ademas de la falla en la sincronizacion, cuando el otro thread finalmente llega y ejecuta `importante()`, al hacer `permiso.wait()` queda bloqueado para siempre, porque ya nadie va a volver a hacer `signal()` sobre esa condition. Es decir, no solo se rompe el orden pedido, sino que se produce un deadlock.

Si siempre el primer thread en entrar al monitor fuera el que ejecuta `importante`, ahí sí se cumpliría el requerimiento — pero el enunciado no garantiza ningún orden de llegada entre los dos threads.

## Punto B 
Si la disciplina utilizada es "signal y continua", tenemos exactamente el mismo problema que el caso anterior. 

Pero tampoco funciona si el orden de ejecucion es primero el thread que hace `importante()` y luego el que hace `antesydespues()`, lo unico que pasa en ese caso es que el que hace `antesydespues()` no va a quedar en deadlock. Pero el orden no se va a respetar porque: 

- El thread `b` entra al monitor y ejecuta `importante()`.
    - Se cuelga en `permiso.wait()`, asi que se bloquea y libera el monitor. 
- El thread `a` entra al monitor y ejecuta `antesydespues()`. 
    - Ejecuta la parte del antes. 
    - Ejecuta `permiso.signal()`, pero como es "signal y continua" no va a abandonar el monitor y cederselo a `b`, si no que va a seguir ejecutando.
    - Ejecuta la parte de despues, incumpliendo el orden. 
    - Abandona el monitor.
- El thread `b` vuelve y ejecuta la parte importante. 

En este caso tenemos una traza valida en la que el orden de ejecución es: `antes->despues->importante`.

# Ejercicio 2 
Primera idea con signal y continue:
```java
monitor SecuenciadorTernario{
    int turnoActual = 1;


    condition esperarTurno;

    void primero(){
        while(turnoActual != 1){
            wait(esperarTurno);
        }
        print("Primero");
        turnoActual = 2;
        signalAll(esperarTurno);
    }

    void segundo(){
        while(turnoActual != 2){
            wait(esperarTurno);
        }
        print("Segundo");
        turnoActual = 3;
        signalAll(esperarTurno);
    }

    void tercero()?{
        while(turnoActual != 3){
            wait(esperarTurno);
        }
        print("Tercero");
        turnoActual = 1;
        signalAll(esperarTurno);
    }
}
```

Una mejor opcion con tres condiciones distintas para no tener que despertar a todos: 
```java
monitor SecuenciadorTernario{
    int turnoActual = 1;


    condition esperarTurno1;
    condition esperarTurno2;
    condition esperarTurno3;

    void primero(){
        while(turnoActual != 1){
            wait(esperarTurno1);
        }
        print("Primero");
        turnoActual = 2;
        signal(esperarTurno2);
    }

    void segundo(){
        while(turnoActual != 2){
            wait(esperarTurno2);
        }
        print("Segundo");
        turnoActual = 3;
        signal(esperarTurno3);
    }

    void tercero(){
        while(turnoActual != 3){
            wait(esperarTurno3);
        }
        print("Tercero");
        turnoActual = 1;
        signal(esperarTurno1);
    }
}
```

Una solucion posible si se usa signal and wait puede ser: 

```java
monitor SecuenciadorTernario{
    int turnoActual = 1;

    condition esperarTurno1;
    condition esperarTurno2;
    condition esperarTurno3;

    void primero(){
        if(turnoActual != 1){
            wait(esperarTurno1);
        }
        print("Primero");
        turnoActual = 2;
        signal(esperarTurno2);
    }

    void segundo(){
        if(turnoActual != 2){
            wait(esperarTurno2);
        }
        print("Segundo");
        turnoActual = 3;
        signal(esperarTurno3);
    }

    void tercero(){
        if(turnoActual != 3){
            wait(esperarTurno3);
        }
        print("Tercero");
        turnoActual = 1;
        signal(esperarTurno1);
    }
}
```

# Ejercicio 3 
Un primer approach puede ser:

- Para signal and wait:
    ```java
    monitor Barrera(int N){
        int contador = 0;

        condition esperar;

        void esperar(){
            if(contador < N){
                contador++; 
                if(contador == N){
                    signalAll(esperar);
                }
                else{
                    wait(esperar);
                }
            }   
        }
    }
    ```
- Para signal and continue: 
    ```java
    monitor Barrera(int N){
        int contador = 0;

        condition esperar;

        void esperar(){
            if(contador < N){
                contador++; 
                if(contador == N){
                    signalAll(esperar);
                }
                else{
                    while(contador != N){
                        wait(esperar);
                    }
                }
            }
            
        }
    }
    ```
- Si quiero que no se haga un `signalAll` si no que se vayan levantando uno por uno una solucion usando signal and continue puede ser: 
    ```java
    monitor Barrera(int N){
        int contador = 0;
        int cantidadPorLevantar = N - 1; // Siempre va a haber uno que no hace falta despertar pq nunca se duerme

        condition esperar;

        void esperar(){
            if(contador < N){
                contador++; 
                if(contador == N){
                    cantidadPorLevantar -= 1;
                    signal(esperar);
                }
                else{
                    while(contador != N){
                        wait(esperar);
                    }
                    if(cantidadPorLevantar > 0){
                        cantidadPorLevantar -= 1; 
                        signal(esperar);
                    }
                }
            }
            
        }
    }
    ```

## Punto B 
- Con signal and wait, usando un `signalAll`: 
    ```java
    monitor Barrera(int N){
        int contador = 0;

        condition esperar;

        void esperar(){
            if(contador < N){
                contador++; 
                if(contador == N){
                    contador = 0;
                    signalAll(esperar);
                }
                else{
                    wait(esperar);
                }
            }   
        }
    }
    ```
- Con signal and continue, si usamos un `signalAll`, ya no nos sirve hacer un `while` pidiendo que el contador sea igual a `N`, porque justamente en el codigo anterior el contador se reinicia, haciendo que se bloqueen todos los procesos. Por lo tanto: 
    ```java
    monitor Barrera(int N){
        int contador = 0;
        int uso = 0;

        condition esperar;

        void esperar(){
            int miUso = uso;
            if(contador < N){
                contador++; 
                if(contador == N){
                    contador = 0;
                    uso += 1;
                    signalAll(esperar);
                }
                else{
                    while(miUso == uso){
                        wait(esperar);
                    }
                }
            }   
        }
    }
    ```
# Ejercicio 4 
## Punto A 
Para que ocurra que `liberar` nunca se bloquee, obligatoriamente hay que usar signal and continue, porque en signal and wait siempre que le demos un signal con alguien esperando, el thread que esta en `liberar` va a dejar que termine el otro thread antes. 

Tambien para esta solución asumo que no pueden haber sporious wakeups. Si los hubiera habría que complejizar el `esperar` haciendo que haga un loop sobre un while, y en ese caso si que puede pasar que alguien que acaba de entrar pase derecho y se saltee a alguien que estaba dormido. 

```java
Monitor Atrapador{
    int cantidadEsperando = 0;

    condition esperar; 

    void esperar(){
        cantidadEsperando++; 
        wait(esperar);
    }

    void liberar(int N){
        if(cantidadEsperando >= N){
            for(int i = 0; i < N; i++){
                signal(esperar);
                cantidadEsperando--;
            }
        }
    }
}
```

Por como funcionan los monitores y las variables de condicion, no puede ocurrir que un thread `T2` "no tenga que esperar", si o si una vez que ejecute el `wait` va a ceder el acceso al monitor y va a tener que esperar a que un thread ejecute `liberar` y haya los suficientes threads para despertar en ese llamado. 

Lo que si puede ocurrir es si por ejemplo `T1` llama a `esperar`, y luego `T2` llama a `esperar`, y despues otro thread los despierta a los dos, por mas que `T1` va a ser despertado antes que `T2` porque las condiciones son colas FIFO, nada le garantiza a cada una el orden en el que van a volver a ejecutar, ya que van a tener que pelearse por el acceso al monitor con el resto de los threads sin un orden. Por lo que `T2` puede "terminar" la ejecucion de `esperar` antes que `T1` por mas que `T1` la inicio antes. 

## Punto B 
Usando signal and continue una primera solución puede ser:
```java
Monitor Atrapador{
    int cantidadEsperando = 0;
    bool hayLiberador = false; 
    int cantidadNecesariaParaLiberar = -1;
    condition esperar; 
    condition puedeLiberar; 
    condition hayProcesosSuficientesParaLiberar;

    void esperar(){
        cantidadEsperando++; 
        if(cantidadEsperando == cantidadNecesariaParaLiberar){
            signal(hayProcesosSuficientesParaLiberar);
        }
        wait(esperar);
    }

    void liberar(int N){
        while(hayLiberador){
            wait(puedeLiberar);
        }
        hayLiberador = True; // Soy yo
        cantidadNecesariaParaLiberar = N;
        while(cantidadEsperando < N){
            wait(hayProcesosSuficientesParaLiberar);
        }
        // Hay la cantidad suficiente para liberar, asi que los liberamos uno por uno. 
        for(int i = 0; i < N; i++){
            signal(esperar);
            cantidadEsperando--;
        }
        hayLiberador = False;
        signal(puedeLiberar);
    }
}
```
Es necesario usar el semaforo para indicar si hay un liberador, porque si no mientras el liberador esta dormido o incluso alguien lo despierta y tiene que volver a esperar para entrar al monitor, en el medio alguien le puede ganar de manos y entrar al monitor y pisarle cuantos esta esperando para liberar, por lo que si vuelve a entrar el liberador original le va a pisar el valor compartido de cuantos estaba esperando para liberar, haciendo que deje de tener sentido la variable compartida. 

Otra solución lo que puede hacer es en vez de hacer que quien es el liberador sea estrictamente FIFO, y si hay algun liberador que ya puede liberar se tenga que quedar esperando, podemos hacer que todos los liberadores se duerman en la misma variable de condición y se tire un `signalAll` en cada paso, si algún liberador al despertarte puede liberar, lo va a hacer.  

```java
Monitor Atrapador{
    int cantidadEsperando = 0;
    condition esperar;  
    condition hayProcesosSuficientesParaLiberar;

    void esperar(){
        cantidadEsperando++; 
        signalAll(hayProcesosSuficientesParaLiberar);
        wait(esperar);
    }

    void liberar(int N){
        while(cantidadEsperando < N){
            wait(hayProcesosSuficientesParaLiberar);
        }
        // Hay la cantidad suficiente para liberar, asi que los liberamos uno por uno. 
        for(int i = 0; i < N; i++){
            signal(esperar);
            cantidadEsperando--;
        }
        signal(puedeLiberar);
    }
}
```

Una solución para no tener que hacer `signalAll` cada vez que algún proceso se va a dormir, independientemente de si hay la cantidad suficiente de procesos dormidos para que al menos un liberador pueda trabajar, lo que podemos hacer es tener un set de cantidades necesarias para que al menos un liberador pueda trabajar, si la cantidad actual pertenece a ese conjunto despertamos a todos, pero sabemos que al menos uno va a poder trabajar.

Para poder usar un conjunto, como puede pasar que dos liberadores distintos quieran liberar a la misma cantidad de procesos, tenemos que en el conjunto en vez de meter numeros enteros tenemos que meter tuplas que la primer componente represente el numero de procesos que tiene que despertar, y el segundo la aparicion de ese elemento dentro del conjunto. 

Si no hacemos esto si dos liberadores distintos quieren liberar a `N` procesos, el primero que lo haga va a despertarlos y a eliminar `N` del conjunto, dejando en deadlock al segundo proceso.

Supongo que tengo implementada una función `perteneceAPrimerComponente(Set<(int, int)> cjto, int e)` que devuelve verdadero si `e` es una primer componente del conjunto, falso si no. 

Supongo que tengo implementada una función `siguienteIndiceDeElemento(Set<(int, int)> cjto, int e)` que devuelve el siguiente indice de un elemento dentro de un conjunto. Por ejemplo si el conjunto es `{(3,0), (2,1), (3,1)}` va a devolver 2, porque el siguiente indice para el 3 libre es el 2. Para implementarlo lo que se podría hacer es un for donde preguntas si `(e, 0)` pertenece, si pertenece preguntas si pertenece `(e,1)` asi hasta que sea falso que encontramos nuestro indice. 

```java
Monitor Atrapador{
    int cantidadEsperando = 0;
    Set<(int, int)> cantidadesNecesarias = {};
    condition esperar;  
    condition hayProcesosSuficientesParaLiberar;

    void esperar(){
        cantidadEsperando++; 
        if(perteneceAPrimerComponente(cantidadEsperando, cantidadesNecesarias)){
            signalAll(hayProcesosSuficientesParaLiberar)
        }
        wait(esperar);
    }

    void liberar(int N){
        int indice = siguienteIndiceDeElemento(cantidadesNecesarias, N);
        cantidadesNecesarias.add((N, indice));
        while(cantidadEsperando < N){
            wait(hayProcesosSuficientesParaLiberar);
        }
        // Hay la cantidad suficiente para liberar, asi que los liberamos uno por uno. 
        for(int i = 0; i < N; i++){
            signal(esperar);
            cantidadEsperando--;
        }
        cantidadesNecesarias.remove((N, indice));
        signal(puedeLiberar);
    }
}
```

TODO: Solucion sin signalAll. 

# Ejercicio 5 
Asumo que no hay una cota de sillas de espera, o de personas que pueden estar esperando y que puede haber mas de un peluquero a la vez cortando. 
```java
Monitor Pelu(){
    condition hayCliente; 
    condition hayPeluquero; 
    condition terminoCorte;

    int cantidadPersonasEsperando = 0;
    int proximoNumeroDeCorte = 0;
    int proximoNumeroAAtender = 0;
    Set<int> numerosLlamados = {};      // turnos ya tomados por algun peluquero
    Set<int> cortesTerminados = {};     // turnos cuyo corte ya termino

    local int numeroClienteSiendoAtendido = -1;

    void cortarseElPelo(){
        cantidadPersonasEsperando += 1;
        int miNumero = proximoNumeroDeCorte; 
        proximoNumeroDeCorte += 1;
        signalAll(hayCliente);

        while(!numerosLlamados.contains(miNumero)){
            wait(hayPeluquero);
        }

        while(!cortesTerminados.contains(miNumero)){
            wait(terminoCorte);
        }
        cortesTerminados.remove(miNumero);   
    }

    int empezarCorte(){
        while(cantidadPersonasEsperando < 1){
            wait(hayCliente);
        }
        cantidadPersonasEsperando -= 1;

        numeroClienteSiendoAtendido = proximoNumeroAAtender;
        proximoNumeroAAtender += 1;
        numerosLlamados.add(numeroClienteSiendoAtendido);
        signalAll(hayPeluquero);   
    }

    void terminarCorte(){
        cortesTerminados.add(numeroClienteSiendoAtendido);
        signalAll(terminoCorte);
    }
}
```

## Punto B
Si hay que usar una sola variable de condicion, se podría reemplazar el código actual y en todas las variables de condición usar una sola y el código seguiria funcionando, ya que siempre que alguien se despierta verifica que se cumpla una condición, así que si alguien se despierta y ve que su condición no se cumple se va a volver a dormir. 

Va a ser mucho mas ineficiente la ejecución porque al hacer un `signalAll` siempre sobre la misma cola van a estar todo el tiempo despertandose procesos y consumiendo computo procesos que siguen sin tener los recursos necesarios desbloqueados. Para evitar esto se podrían implementar colas de semaforos para cosas como por ejemplo el manejo de cuando un peluquero le esta cortando al cliente el cliente se duerma en ese semaforo en particular y que el peluquero lo despierte a el solo al terminar el corte. 

# Ejercicio 6
## Punto A
Asumo por lo que dice la consigna que siempre va a haber un solo orador que va repitiendo la charla. 

La idea del monitor es que van a haber cuatro metodos, para los asistentes `entrarASala`, `salirSala`, `comenzarCharla` y `terminarCharla`. 

La idea va a ser que: 
- En `entrarSala` si hay alguien hablando la persona se va a quedar esperando afuera hasta que la charla termine, y si la charla termino pero todavia hay 50 personas dentro, va a tener que seguir esperando. Una vez que entramos vamos a tener que marcar que somos una persona mas dentro incrementando la variable compartida que dice la cantidad de gente que hay.
- En `salirSala` si hay alguien hablando vamos a tener que esperar hasta que termine, y una vez podamos salir antes de salir vamos a decrementar en uno la cantidad de gente dentro de la sala. 
- En `comenzarCharla` va a primero que verificar que haya al menos una persona, si no lo hay descansa 5 minutos mientras el resto de la gente puede entrar, por lo que hay que liberar el acceso al monitor. 
- En `finalizarCharla` va a tener que obligatoriamente descansar 5 minutos, y mientras tanto liberar el monitor para que pueda entrar y salir gente. 

En esta primer versión asumo que una persona puede entrar a una charla y salir sin que arranque la charla: 

```java 
Monitor Conferencia(){
    int cantidadAsistentes = 0;
    boolean hayCharla = false;

    condition hayLugar;
    condition terminoCharla; 

    void entrarSala(){
        while(true){
            if(hayCharla){
                wait(terminoCharla);
            }
            else if(cantidadAsistentes >= 50){
                wait(hayLugar); // Si alguien se despierta de aca, va a tener que verificar primero que no haya charla y que haya lugar para salir del ciclo.
            }
            else{
                break;
            }
        }
        // Si estoy aca es pq entre 
        cantidadAsistentes += 1;
    }

    void salirSala(){
        while(hayCharla){
            wait(terminoCharla);
        }
        // Si estoy aca es pq sali 
        cantidadAsistentes -= 1;
        signalAll(hayLugar);
    }

    boolean comenzarCharla(){
        if(cantidadAsistentes == 0){
            return false;
        }
        // Hay gente asi que cerramos la puerta 
        hayCharla = true; 
        return true; 
    }

    void finalizarCharla(){
        hayCharla = false; 
        signalAll(terminoCharla);
        // Dormir 5 min y liberar el acceso al monitor 
    }
}
```
Y el thread orador va a tener que correr un ciclo: 
```java
Thread orador(Conferencia conferencia) {
    while (true) {
        boolean comenzo = conferencia.comenzarCharla();

        if (comenzo) {
            darCharla(); // fuera del monitor
            conferencia.finalizarCharla();
        }

        Thread.sleep(5 * 60 * 1000); // fuera del monitor
    }
}
```

Para que una persona no se pueda ir entre que entra la sala y que termina la charla, con una variable local de la persona que sea el numero de charla actual, y tengo otra variable global que es el numero de charla, y que para salir tiene que hacer un wait a una variable de condicion encerrada en un ciclo que pida que el numero de charla por empezar tiene que ser estrictamente mayor al numero de charla en la que entro por primera vez.  

## Punto B
Si ahora hay 3 personas que se van a "pelear" por dar la charla, pero no van a hacer mas los oradores los 5 minutos de descanso al finalizar la charla o esperar 5 minutos si no hay nadie, si no que una vez que logren entrar al auditorio van a esperar a que hayan 40 personas y ahí van a arrancar la charla. 

```java 
Monitor Conferencia(){
    static final int CAPACIDAD_SALA = 50;
    static final int MIN_PARA_ARRANCAR = 40;

    int cantidadAsistentes = 0;
    boolean hayCharla = false;
    boolean hayOradorPresente = false;

    condition hayLugar;
    condition terminoCharla;
    condition hayOrador;
    condition haySuficienteGenteParaEmpezar;

    void entrarSala(){
        while(true){
            if(hayCharla){
                wait(terminoCharla);
            }
            else if(cantidadAsistentes >= CAPACIDAD_SALA){
                wait(hayLugar); // Al despertar, vuelve a verificar que no haya charla y que haya lugar antes de salir del ciclo.
            }
            else{
                break;
            }
        }
        // Si estoy aca es pq entre
        cantidadAsistentes += 1;
        if(cantidadAsistentes == MIN_PARA_ARRANCAR){
            signal(haySuficienteGenteParaEmpezar);
        }
    }

    void salirSala(){
        while(hayCharla){
            wait(terminoCharla);
        }
        // Si estoy aca es pq sali
        cantidadAsistentes -= 1;
        signalAll(hayLugar);
    }

    void comenzarCharla(){
        while(hayOradorPresente){
            wait(hayOrador);
        }
        // Si estamos aca es pq ganamos la carrera y somos el orador actual
        hayOradorPresente = true;
        while(cantidadAsistentes < MIN_PARA_ARRANCAR){
            wait(haySuficienteGenteParaEmpezar);
        }
        // Hay la cantidad de gente deseada asi que cerramos la puerta
        hayCharla = true;
    }

    void finalizarCharla(){
        hayCharla = false;
        signalAll(terminoCharla);
        hayOradorPresente = false;
        signalAll(hayOrador);
    }
}
```
Y el thread orador va a tener que correr un ciclo: 
```java
Thread orador(Conferencia conferencia) {
    while (true) {
        conferencia.comenzarCharla();
        darCharla(); // fuera del monitor
        conferencia.finalizarCharla();
    }
}
```

Para evitar que una persona se quede encerrada dentro de las charlas infinitamente o que siempre gane un orador y de la charla siempre el se podría implementar un turnstile para ordenar quienes entran al semaforo, de forma que por ejemplo una persona que quiere irse de la charla tiene que esperar a lo sumo a 3 charlas para irse en un peor caso. 


# Ejercicio 7 
## Punto A
```java 
Thread Jugador(monitor juego){
    boolean gane = false;
    string palabra; 
    string monto;
    while(!juego.concluido()){
        palabra = Str.random();
        monto = Int.random();
        gane = juego.apostar(palabra, monto);
    }
    if(gane){
        System.out.println("Gane con la palabra: " + palabra + ", gané " + (monto * 10) + " pesos");
    }
    else{
        System.out.println("Perdi");
    }
}
```

## Punto B 
```java 
monitor Juego(String palabra) {
    boolean alguienAdivino = false;
    Thread ultimoJugador = null;

    condition puedeJugar;

    boolean apostar(String palabraPropuesta, int apuesta) {
        Thread yo = Thread.currentThread();

        // Mientras el último que jugó sea yo mismo, espero que juegue otro.
        while (yo == ultimoJugador) {
            wait(puedeJugar);
        }

        if (alguienAdivino) {
            ultimoJugador = yo;
            return false;
        }

        boolean gane = palabraPropuesta.equals(palabra);
        if (gane) {
            alguienAdivino = true;
        }

        ultimoJugador = yo;
        signalAll(puedeJugar);
        return gane;
    }

    boolean concluido() {
        return alguienAdivino;
    }
}
```