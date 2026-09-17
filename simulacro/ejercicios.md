# Ejercicio 1 
## Punto A 
Los requerimientos de exclusión mutua son: 
- Exclusión Mutua (Safety): En cualquier momento hay como máximo un
proceso en la región crítica.
- Ausencia de deadlocks (Liveness): Si varios procesos intentan entrar a la
sección crítica tarde o temprano alguno entra
- Ausencia de inhanición (garantía de entrada) (Liveness): Cualquier
proceso que intenta entrar a su sección crítica tarde o temprano entra.

El problema que va a haber y lo que no se cumple es la propiedad de ausencia de deadlocks.

El caso de deadlock puede ocurrir si: 

Supongamos que tenemos tres procesos corriendo concurrentemente el código, y separamos la función de anotarse en los pasos atomicos, donde cada `pn` es un paso atomico 
```java 
bool anotarse(int idThread){
    p1: esperando[idThread] = true; 
    p2: local esperando = cantEsperando; 
    p3: esperando = esperando + 1; 
    p4: cantEsperando = esperando;
    p5: return (cantEsperando == 1);
}
```

El problema que puede ocurrir, es que si dos threads concurrentemente ejecutan `anotarse`, lo que puede ocurrir es que el primero ejecute hasta `p4` y lo corten antes de ejecutar `p5`, por lo que en ese punto `cantEsperando` va a ser igual a 1. Si viene otro proceso y ejecuta `anotarse`, este va a incrementar `esperando` a 2, haciendo que `return (cantEsperando == 1)` siempre sea falso, de forma que ningún thread va a poder ejecutar `llamarProximo()`, haciendo que caigan en deadlock todos los procesos ya que ninguno va a siquiera poder entrar a la sección critica. 

Para romper exclusion mutua tenes 3 threads, el mas chico espera, los otros dos hacen dos anotarse. 

## Punto B 
En caso de que las tres funciones sean atomicas, si se cumple exclusión mutua, justificando cada propiedad por separado: 
- Se cumple exclusión mutua porque, si varios threads ejecutan concurrentemente la primer parte antes de entrar a la sección critica, solamente a un thread le va a dar verdadero soy primero, de forma que `llamarProximo` solo se va a ejecutar una sola vez. Por lo que `proximo` va a tener un `id` que es igual al de alguno de los threads. Por lo que uno de los threads va a poder entrar a la sección critica. En caso de que se haga un context switch al thread dentro de la sección critica, nadie va a poder hacer un `llamarProximo` antes de la entrada a la sección critica. 
    Tarde o temprano el thread que entro a la sección critica va a llamar a `soyUltimo`, y `llamarProximo`, ese `llamarProximo` y otro va a ser el siguiente, por lo que siempre va a haber uno solo dentro de la sección critica. 
- Se cumple ausencia de deadlocks porque el unico lugar que puede bloquear a los procesos es el ciclo `while(proximo != id)`. Como siempre `proximo` va a tener el valor del id de un thread porque `anotarse` nunca puede devolver falso al ser atomico, siempre se van a ir asignando ids y tarde o temprano el `id` de cada thread que se anote va a ser llamado. 
- Como el `N` es fijo, y no puede ocurrir que todo el tiempo entren procesos a rolete, y se van a ir pasando el acceso cada uno uno a la vez, tarde o temprano incluso el que tenga id mayor va a poder entrar. En caso de que esten dentro de un `while true` los procesos que tengan id mayor podrían tener inanición al siempre darle el paso a los procesos de id mayor. 

Para probar mutex poner como absurdo que hay dos threads en la seccion critica. Si eso paso es porque uno de los dos entra primero y sale primero del while, si uno salio primero del while la condicion del otro si tambien entro pudo salir del while, pero eso no puede pasar pq `proximo != id`, como `id` es unico, la unica manera de que cambie `id`, es que alguien dentro de la seccion critica lo cambie, o que alguien de afuera sea que llamar proximo venga de afuera. Como las funciones son atomicas lo puede hacer solo alguien que acaba de llegar o alguien que esta saliendo (que no pase nada), y un llegar proximo por el que alguien llega no puede pasar porque cantidad esperando esta protegido, si hay threads que estan esperando, ese numero va a ser 0, eventualmente va a ser 1 y nunca puede pasar que decremente con esa persona este dentro. Entonces es estrictamnete mayor a 1, luego nadie va a llamar a llamarProximo hasta que salga de la sección critica. 

Para ausencia de deadlock mostrar que ninguno pueda quedarse colgado en el while, separarlo en que si o si cuando alguien llega proximo no va a ser igual a -1, y que se va a ir cambiando por cadena. 

Para ausencia de inanhicion decir x absurdo que supongamos que nunca ejecuta, entonces tiene que ser pq siempre llegan nuevos, pero eso no puede pasar. 

## Punto C 
El problema es que pasa si estan ejecutando en mas de un core y vos tengas un dato en una cache, y otro nucleo lo pisa. 

El problema es con las variables `volatile`

```
volatile D
D++ 
b
```

Siempre `D++` se va a ejecutar antes que `b`. Aca el otro core 

Dos threads ven cant esperando sea igual a 1, no pq se interfirieron, pero pq no ven el cambio hecho remotamente. 

El cambio deberia ser cambiar las variables problematicas compartidas a `volatile`. 

# Ejercicio 2 
El problema es una combinacion de filosofos comensales + lectores y escritores. 

La idea va a ser que cada cuenta tenga un mutex asociado a su dato, y tambien tenemos un mutex para proteger una variable de cantidadDeTransferencias, la idea es hacer un cambio de acceso cuando no hay transferencias en curso y en ese caso se puede hacer el reporte. 

```java
global semaphore[] mutexs = [Semaphore(1), ...., Semahpore(1)]
global int[] balances = [0, ...., 0];
global cantidadTransferencias = 0;
global semaphore mutex = Semaphore(1);
global semaphore molinete = Semaphore(1, True);


boolean transferir(int origen, int destino, int monto){
    molinete.acquire();
    molinete.release();
    mutex.acquire();
    if(cantidadTransferencias == 0){
        acceso.acquire();
    }
    cantidadTransferencias += 1;
    mutex.signal();
    mutexs[min(origen, destino)].acquire();
    mutexs[max(origen, destino)].acquire();
    // if para ver si son negativos o no
    balances[origen] -= monto; 
    balances[destino] += monto;
    mutexs[max(origen, destino)].release();
    mutexs[min(origen, destino)].release(); 
    mutex.acquire();
    cantidadTransferencias -= 1;
    if(cantidadTransferencias == 0){
        acceso.release();
    }
    mutex.signal();
}
```

```java 
void cuentasEnRojo(){
    molinete.acquire();
    acceso.acquire();
    for(...)
    acceso.release();
    molinete.release();
}
```

Siempre como patron desbloquear en orden inverso al que tomar. 

## Punto B 
Haces que todos los semaforos del array sean fuertes y el acceso tambien. 

## Punto C
Tomas la lista de destino y la ordenas, y tomas todos los locks de los destinos en orden y dp haces las transferencias. 

```java 
cuentas.sort();
for (c in cuentas){
    mutexs[c].adquire();
}

// Critica 

for (c in cuentas.reverse()){
    mutexs[c].release();
}
```

# Ejercicio 3 
```java 
Monitor NDJ(){
    condition lugarDisponible;

    void sentarse(){
        while(!existe(i) tal que libre[i] y lugares[i] >= k){
            wait(lugarDisponible)
        }
        i = aquella mesa que sea libre y tenga lugar 
        libre[i] = false
        return i;
    }

    void levantarse(int mesa){

    }
}
```

Para que no haya starvation habria que usar un numero de ticket, una solucion que lo secuencializa bastante es que al while le agregues que el proximo no sea igual al ticket, pero lo sequencializa bastante, para mejorarlo se podria usar un conjunto de tuplas `(cantidadGente, ticket)` y cuando ves el while gana el que entra en la mesa que menor numero de ticket tiene. 

Otra solución posible es con colas. Tenemos una lista donde tenemos una tupla cual es mi ticket y una variable de condicion. Cuando llega alguien se agrega a la lista, cuando alguien se despierta puede buscar de todos los pedidso cual es el mas cercano que puede entrar a la mesa k, y directamente levantas a ese grupo que queda con una variable de condicion dinamica. Tenes una variable de condicion por cada pedido. Levantas al primero que cumple de la cola/ lista que le puede venir bien la mesa que liberaste. 

## Punto B 
Habria que cambiar los whiles por if. 

# Ejercicio 4 
## Punto A

```java 
lock[] locks = LOCK[n];
void inc(int id){
    locks[id % tamano].lock();
    valor[id % tamano]++;
    locks[id % tamano].unlock();
}

int get(){
    int total = 0;
    for(int i = 0; i < n; i++){
        total = valor[i];
    }
    return total;
}
```
Hay que tomar todos los locks en el get porque si no hay problemas de linearabizilidad porque podes estar sumando y que en el medio te pueden sumar algo de mas atras.

En el incremento el punto de linealizacion es cuando haces el cambio que sumas en uno. 

Los `locks` son fair

## Punto C
```java 
void inc(int id){
    bucket = id % n;
    actual = valor[bucket];
    while(true){
        if(cas(valor[bucket], actual, actual + 1)){
            return;
        }
    }
}
```

## Punto D
Es linealizable, el punto de linealizacion es cuando terminas el segundo loop por ultima vez. 

## Punto E 
No es wait free porque si siempre estan metiendose se queda colgado. 

## Inciso extra 
Como el cas no tiene el lock los incrementos los ignora, pasa que la linealizacion se romperia si trataria.