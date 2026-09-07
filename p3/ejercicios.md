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
```
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
```
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

```
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
    ```
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
    ```
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
    ```
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
    ```
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
    ```
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

```
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
Usando signal and continue: 
```
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

Se podria permitir algo asi y que cada vez que alguien libera mande un signal a `hayProcesosSuficientesParaLiberar` y que cada liberador solo verifique contra si la cantidad esperando es igual a su parametro, pero esto agregaria muchos signals de mas y podria hacer que algun liberador que esta en condiciones de liberar no libere porque otro le robo el lugar en la condicion o esta afuera del monitor esperando a entrar de nuevo. 

# Ejercicio 5 
Asumo que no hay una cota de sillas de espera, o de personas que pueden estar esperando y que puede haber mas de un peluquero a la vez cortando. 
```
Monitor Pelu(){
    condition hayCliente; 
    condition hayPeluquero; 
    condition terminoCorte;

    int cantidadPersonasEsperando = 0;
    int proximoNumeroDeCorte = 0;
    int proximoNumeroAAtender = 0;
    Set<int> numerosLlamados = {};      // turnos ya tomados por algun peluquero
    Set<int> cortesTerminados = {};     // turnos cuyo corte ya termino

    void cortarseElPelo(){
        cantidadPersonasEsperando += 1;
        int numeroClienteActualSiendoAtendido = proximoNumeroDeCorte; 
        proximoNumeroDeCorte += 1;
        signalAll(hayCliente);

        while(!numerosLlamados.contains(numeroClienteActualSiendoAtendido)){
            wait(hayPeluquero);
        }

        while(!cortesTerminados.contains(numeroClienteActualSiendoAtendido)){
            wait(terminoCorte);
        }
        cortesTerminados.remove(numeroClienteActualSiendoAtendido);   
    }

    int empezarCorte(){
        while(cantidadPersonasEsperando < 1){
            wait(hayCliente);
        }
        cantidadPersonasEsperando -= 1;

        int numeroClienteActualSiendoAtendido = proximoNumeroAAtender;
        proximoNumeroAAtender += 1;
        numerosLlamados.add(numeroClienteActualSiendoAtendido);
        signalAll(hayPeluquero);

        return numeroClienteActualSiendoAtendido;   
    }

    void terminarCorte(int numeroClienteActualSiendoAtendido){
        cortesTerminados.add(numeroClienteActualSiendoAtendido);
        signalAll(terminoCorte);
    }
}
```

## Punto B
Si hay que usar una sola variable de condicion, todos podrian usar la misma, como todos por mas que se despiertan luego ven que valga el booleano, no va a ocurrir nunca que alguien se saltee una condicion. Si fuera un stop and wait sin while, ahi si tendriamos problemas y seria necesario volver a verificar cada condicion luego de despertarnos. 

