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

Supongamos que tenemos tres procesos corriendo concurrentemente el código, y separamos la función de anotarse en los pasos atómicos, donde cada `pn` es un paso atómico 
```java 
bool anotarse(int idThread){
    p1: esperando[idThread] = true; 
    p2: local esperando = cantEsperando; 
    p3: esperando = esperando + 1; 
    p4: cantEsperando = esperando;
    p5: return (cantEsperando == 1);
}
```

El problema que puede ocurrir, es que si dos threads concurrentemente ejecutan `anotarse`, lo que puede ocurrir es que el primero ejecute hasta `p4` y lo corten antes de ejecutar `p5`, por lo que en ese punto `cantEsperando` va a ser igual a 1. Si viene otro proceso y ejecuta `anotarse`, este va a incrementar `esperando` a 2, haciendo que `return (cantEsperando == 1)` siempre sea falso, de forma que ningún thread va a poder ejecutar `llamarProximo()`, haciendo que caigan en deadlock todos los procesos ya que ninguno va a siquiera poder entrar a la sección crítica al proximo siempre valer -1 y los ids de los thread siempre valen distinto a -1. 

Para romper exclusión mutua tenés 3 threads, el más chico espera, los otros dos hacen dos anotarse (TODO traza).

## Punto B 
En caso de que las tres funciones sean atómicas, sí se cumple exclusión mutua, justificando cada propiedad por separado: 
- Se cumple exclusión mutua porque, si varios threads ejecutan concurrentemente la primer parte antes de entrar a la sección crítica, solamente a un thread le va a dar verdadero soy primero, de forma que `llamarProximo` solo se va a ejecutar una sola vez. Por lo que `proximo` va a tener un `id` que es igual al de alguno de los threads. Por lo que uno de los threads va a poder entrar a la sección crítica. En caso de que se haga un context switch al thread dentro de la sección crítica, nadie va a poder hacer un `llamarProximo` antes de la entrada a la sección crítica. 
    Tarde o temprano el thread que entró a la sección crítica va a llamar a `soyUltimo`, y `llamarProximo`, ese `llamarProximo` y otro va a ser el siguiente, por lo que siempre va a haber uno solo dentro de la sección crítica. 
- Se cumple ausencia de deadlocks porque el único lugar que puede bloquear a los procesos es el ciclo `while(proximo != id)`. Como siempre `proximo` va a tener el valor del id de un thread porque `anotarse` nunca puede devolver falso al ser atómico, siempre se van a ir asignando ids y tarde o temprano el `id` de cada thread que se anote va a ser llamado. 
- Como el `N` es fijo, y no puede ocurrir que todo el tiempo entren procesos a rolete, y se van a ir pasando el acceso cada uno uno a la vez, tarde o temprano incluso el que tenga id mayor va a poder entrar. En caso de que estén dentro de un `while true` los procesos que tengan id mayor podrían tener inanición al siempre darle el paso a los procesos de id mayor. 

Para probar mutex poner como absurdo que hay dos threads en la sección crítica. Si eso pasó es porque uno de los dos entra primero y sale primero del while, si uno salió primero del while la condición del otro si también entró pudo salir del while, pero eso no puede pasar pq `proximo != id`, como `id` es único, la única manera de que cambie `id`, es que alguien dentro de la sección crítica lo cambie, o que alguien de afuera sea que llamar próximo venga de afuera. Como las funciones son atómicas lo puede hacer solo alguien que acaba de llegar o alguien que está saliendo (que no pase nada), y un llegar próximo por el que alguien llega no puede pasar porque cantidad esperando está protegido, si hay threads que están esperando, ese número va a ser 0, eventualmente va a ser 1 y nunca puede pasar que decremente con esa persona esté dentro. Entonces es estrictamnete mayor a 1, luego nadie va a llamar a llamarProximo hasta que salga de la sección crítica. 

Para ausencia de deadlock mostrar que ninguno pueda quedarse colgado en el while, separarlo en que sí o sí cuando alguien llega próximo no va a ser igual a -1, y que se va a ir cambiando por cadena. 

Para ausencia de inanhición decir x absurdo que supongamos que nunca ejecuta, entonces tiene que ser pq siempre llegan nuevos, pero eso no puede pasar ninguno se va a quedar sin ejecutar. Si estuvieran dentro de un while true ahí si depende del scheduler. 

## Punto C 
El problema es que pasa si están ejecutando en más de un core y vos tengas un dato en una caché, y otro núcleo lo pisa. 

El problema es con las variables `volatile`

```
volatile D
D++ 
b
```

Siempre `D++` se va a ejecutar antes que `b`. Acá el otro core 

Dos threads ven cant esperando sea igual a 1, no pq se interfirieron, pero pq no ven el cambio hecho remotamente. 

El cambio debería ser cambiar las variables problemáticas compartidas a `volatile`, como `proximo` o `cantidadEsperando`. 

# Ejercicio 2 
El problema es una combinación de filósofos comensales + lectores y escritores. 

Las cuentas van a tener ids del `0` al `n`, donde cada nueva cuenta se le asigna un nuevo id y un nuevo balance igual a cero, no incluyo esa parte en la solución. 

La idea va a ser que cada cuenta tenga un mutex asociado a su dato, y también tenemos un mutex para proteger una variable de cantidadDeTransferencias, la idea es hacer un cambio de acceso cuando no hay transferencias en curso y en ese caso se puede hacer el reporte. 

```java
global semaphore[] mutexs = [Semaphore(1), ...., Semahpore(1)]
global int[] balances = [0, ...., 0];
global cantidadTransferencias = 0;
global semaphore mutex = Semaphore(1);
global semaphore molinete = Semaphore(1, True);

global L;

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
    if(balances[origen] - monto < -L){
        mutexs[max(origen, destino)].release();
        mutexs[min(origen, destino)].release(); 
        return false;
    }
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
int cuentasEnRojo(){
    int res = 0;
    molinete.acquire();
    acceso.acquire();
    // Asumo que no puede cambiar el largo de balances en este momento
    for(int i = 0; i < balances.length; i++){
        // Como cualquier transferencia que quiera ocurrir tiene que tener acceso, el cual esta tomado por el thread ejecutando este metodo, no hace falta tomar ningun mutex para consultar el balance de la cuenta, ya que no puede cambiar. 
        if(balances[i] < 0){
            res += 1;
        }
    }
    acceso.release();
    molinete.release();
    return res;
}
```
Como los locks se toman siempre con el mismo orden, no puede ocurrir un deadlock, ya que se elimina el caso cruzado de que un thread 1 quiere procesar una transferencia con origen 1 y destino 2 produsca un deadlock con una que tiene a origen a 2 y destino 1, ya que ambas van a tomar primero el 1 y luego el 2, entonces el que tome primero el 1 gana y no tiene que esperar para tomar el 2, ya que va a estar libre. 

## Punto B 
En la solución del punto anterior puede ocurrir que un pedido que ocurrio antes se procese despues que uno que entro despues, ya que los semaforos en `mutexs` son debiles, de modo que un thread que se queda esperando para modificar un balance puede tener "mala suerte" y que se tarde en darle paso a poder modificar la variable del balance de esa cuenta si todo el tiempo otras transferencias tienen como origen o destino a esa cuenta. 

No hay un limite de cuantas veces puede pasar que una transferencia que se empezo antes tenga que esperar a transferencias posteriores que actuan sobre esa cuenta, porque depende de las transferencias que se quieran hacer y de como el scheduler hace los interleavings. 

Para evitar eso, lo que se puede hacer es que los semaforos de `mutexs` sean fuertes, de forma que siempre la primer transferencia que quiera hacer un movimiento con la cuenta `i` sea la primera en hacerlo. 

Hacés que todos los semáforos del array sean fuertes y el acceso también. 

## Punto C
Tomás la lista de destino y la ordenas, y tomás todos los locks de los destinos en orden y dp hacés las transferencias. 

No puede haber deadlock al tomar los locks en orden. 

```java 
boolean transferirMultiple(int origen, int[] destinos, int[] montos){
    int totalATransferior = 0;
    // Podemos para no tomar todos los locks al pedo primero tomar el del origen y ver si tiene la plata suficiente
    mutexs[origen].acquire();
    for(int i = 0; i < montos.length; i++){
        totalATransferir += montos[i];
    }
    if(totalATransferir > balances[origen] + L){
        mutexs[origen].release();
        return false;
    }
    mutexs[origen].release();
    // Asumo que tengo una función que sortea tanto destinos, como los montos que les debería ir a cada uno para que queden con sentido, y asumo que lo hace por referencia
    // Tambien me devuelve en orden el origen agregado en destinos para tomar los locks en orden y evitar un deadlock.
    destinosMasOrigenSorteados = sortEpico(destinos, montos);
    for(int i = 0; i < destinosMasOrigenSorteados.length; i++){
        mutexs[i].acquire();
    }
    mutexs[origen].acquire();
    for(int i = 0; i < destinos.length; i++){
        balance[origen] -= montos[i];
        balance[destinos[i]] += montos[i];
    }

    for(int i = destinosMasOrigenSorteados.length - 1; i >= 0; i--){
        mutexs[i].release();
    }
    return true; 
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

Para que no haya starvation habría que usar un número de ticket, una solución que lo secuencializa bastante es que al while le agregues que el próximo no sea igual al ticket, pero lo sequencializa bastante, para mejorarlo se podría usar un conjunto de tuplas `(cantidadGente, ticket)` y cuando ves el while gana el que entra en la mesa que menor número de ticket tiene. 

Otra solución posible es con colas. Tenemos una lista donde tenemos una tupla cuál es mi ticket y una variable de condición. Cuando llega alguien se agrega a la lista, cuando alguien se despierta puede buscar de todos los pedidso cuál es el más cercano que puede entrar a la mesa k, y directamente levantás a ese grupo que queda con una variable de condición dinámica. Tenés una variable de condición por cada pedido. Levantás al primero que cumple de la cola/ lista que le puede venir bien la mesa que liberaste. 

## Punto B 
Habría que cambiar los whiles por if. 

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
Hay que tomar todos los locks en el get porque si no hay problemas de linearabizilidad porque podés estar sumando y que en el medio te pueden sumar algo de más atrás.

En el incremento el punto de linealización es cuando hacés el cambio que sumas en uno. 

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
Es linealizable, el punto de linealización es cuando terminás el segundo loop por última vez. 

## Punto E 
No es wait free porque si siempre están metiéndose se queda colgado. 

## Inciso extra 
Como el cas no tiene el lock los incrementos los ignora, pasa que la linealización se rompería si trataría.