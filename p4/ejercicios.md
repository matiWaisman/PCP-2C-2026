# Ejercicio 1 

- En granularidad gruesa: 
    - En el `add` si ocurre que un elemento `b` que no pertenece al conjunto tiene el mismo hash que un elemento `a` que si pertenece, cuando en el código se hace `if (key == curr.key) return false`, la comparación ahora debería ser `if (key == curr.key && curr.value == o.value) return false`. 
    - En el `remove` no habría que cambiar nada porque la desición de eliminar es con un if por el valor del elemento. 
- En granularidad fina:
    - En el `remove` no basta con cortar el ciclo cuando encontramos una clave mayor o igual, porque ahora puede haber varios nodos con la misma `key` y distinto `value`. La condición del ciclo debería ser `while(curr.key < key || (curr.key == key && !curr.item.equals(item)))`, así seguimos avanzando mientras la clave sea menor, o mientras la clave sea igual pero el nodo no sea el elemento que buscamos. Una vez afuera del cuerpo del `while` puede ser por dos cosas: o porque encontramos un elemento cuya key es más grande (o sea que el elemento que queremos eliminar no está), o porque encontramos uno con la misma key que además es `equals()` al que buscamos. Por eso el `if` para eliminarlo cambia a `if(curr.key == key && curr.item.equals(item))`: si es ese elemento lo eliminamos, y si no devolvemos falso.
    - En el `add` la idea sería agregar los elementos en orden, primero ordenando por `key` y después por `value` (dentro de los nodos que comparten `key`). Así que la condición del ciclo debería ser `while(curr.key < key || (curr.key == key && curr.value < o.value))`, avanzando mientras la clave sea menor, o mientras sea igual pero el valor todavía sea menor al que queremos insertar. Al salir del `while` puede ser porque llegamos a un nodo con clave más grande (no hay ningún elemento con esa key todavía, o ya pasamos el lugar que le corresponde por value), o porque encontramos uno con la misma key y mismo value o mayor. Por eso la condición para determinar si el elemento ya pertenece a la lista debería ser `if(curr.key == key && curr.value == o.value)`, devolviendo falso si ya está, e insertando el nuevo nodo entre `pred` y `curr` en caso contrario.
- En la version optimista: 
    - En `validate` no hay que cambiar nada porque en un peor caso donde dos elementos tengan la misma clave van a seguir ambos dentro del ciclo. 
    - En `remove` hay que cambiar la condicion del ciclo `while (curr.key < key)` igual a la version de granularidad fina, `while(curr.key < key || (curr.key == key && curr.value < o.value))` y en el `if(curr.key == key)` habria que agregarle: `if(curr.key == key && curr.item.equals(item))`.
    - En el `add` y `contains` habria que hacer los mismos cambios. 
- En la version lazy: 
    - En el `remove` hay que hacer el mismo cambio al primer ciclo y al `if(curr.key == key)`. 
    - Lo mismo en el `add` y tambien hay que modificar el ciclo del `contains`. 
- En la version lock free: 
    - En el `add` habria que en `if (curr.key == key)` habria que agregar: `if (curr.key == key && curr.item.equals(item))`.

# Ejercicio 2 
Lo que se podria hacer seria una solucion intermedia entre granularidad gruesa y fina en la que tenemos un struct `posicion` donde dentro tiene a las variables `x` e `y`, sumado a un `lock` y otro struct `tamano` que tiene a las variables `alto` y `ancho` sumados a otro `lock`. Asi que en vez de los metodos tener que ser synchronized no hace falta que lo sean, lo primero que se hace al entrar al metodo es intentar de tomar el lock y al terminar se libera. 

No tiene que ser un lock por variable porque en ese caso el punto de linealizacion queda raro, ya que puede ocurrir que alguien llame a `ajustarPosicion`, primero cambia `x`, luego cuando esta esperando por cambiar a `y` otro thread cambia a `x`, por lo que al salir de la funcion el primer thread no va a haber ningun momento en el que la posicion completa esta donde el quiere. 

En este caso el punto de linealizacion va a ser luego de modificar la segunda variable y antes de liberar el lock. 

Si se agregaran nuevas funciones que requieran modificar tanto la posicion como el tamaño lo que habria que hacer es siempre tomar y liberar los locks en el mismo orden, porque si no puede haber una espera circular que hace que terminen en deadlock los threads. 

```java 
public class Posicion {
    int x = 0;
    int y = 0;
    public void lock() { ... }
    public void unlock() { ... }
}

public class Tamano {
    int alto = 0;
    int ancho = 0;
    public void lock() { ... }
    public void unlock() { ... }
}

class Figura {
    Posicion posicion = new Posicion();
    Tamano tamano = new Tamano();

    public void ajustarPosicion() {
        posicion.lock();
        posicion.x = algunX();
        posicion.y = algunY();
        posicion.unlock();
    }

    public void ajustarTamano() {
        tamano.lock();
        tamano.alto = algunAlto();
        tamano.ancho = algunAncho();
        tamano.unlock();
    }
}
```
# Ejercicio 3 
## Punto A 
La idea va a ser que en `assign` antes de darle valor a la posicion del array vamos a bloquear esa posicion para que nadie la pueda swapear y luego de cambiar su valor la desbloqueamos. 

La idea en el `swap` va a ser agarrar los locks de las dos posiciones para que no se les pueda cambiar el contenido durante el swap. Para evitar deadlocks siempre se tiene que agarrar el lock del elemento mas chico y luego el del mas grandes. Si no si se agarra por ejemplo siempre primero el lock de `i` y luego de `j` puede haber un deadlock facil si se hace `swap(1, 2)` se toma el lock del 1 y se hace un context switch a `swap(2,1)` donde se toma el lock de 2 primero y al intentar agarrar el de 1 ocurre un deadlock. 

Los metodos dejarian de pertenecer a un monitor asi que ya se pueden ejecutar concurrentemente. 

```java 
public class Recursos {
    Object[] recursos;
    Lock[] locks;
    int capacidad;

    public Recursos(int n) {
        capacidad = n;
        recursos = new Object[capacidad];
        locks = new Lock[capacidad];
        for (int k = 0; k < capacidad; k++) {
            locks[k] = new Lock();
        }
    }

    public void assign(int pos, Object o) {
        if (pos < capacidad) {
            locks[pos].lock();
            recursos[pos] = o;
            locks[pos].unlock();
        }
    }

    public void swap(int i, int j) {
        if (i < capacidad && j < capacidad) {
            if (i < j) {
                locks[i].lock();
                locks[j].lock();
            } else {
                locks[j].lock();
                locks[i].lock();
            }
            Object aux = recursos[i];
            recursos[i] = recursos[j];
            recursos[j] = aux;
            if (i < j) {
                locks[j].unlock();
                locks[i].unlock();
            } else {
                locks[i].unlock();
                locks[j].unlock();
            }
        }
    }
}
```

Los puntos de linealizacion serian despues de `recursos[pos] = o;` en `assign` y despues de `recursos[j] = aux;` en `swap`. 

## Punto B 
Si los locks son fair, no va a haber inanicion porque siempre el que trate de lockear primero va a tener acceso al lock primero. 

# Ejercicio 4
## Punto A
Un caso en sincronizacion optimista en el que un thread intenta indefinidamente eliminar a un nodo sin conseguirlo puede ocurrir si siempre llega a el `curr` y `pred` adecuados, pero antes de lockearlos alguien agrega un elemento nuevo antes de `curr` reemplazando a `pred`.

El problema es que cuando los lockee, ya no va a seguir valiendo que `pred` apunta a `curr`, si no que hay un nuevo `pred`, por lo que el `validate` va a fallar y va a volver a tener que recorrer los elementos de la lista. Para que haya inanicion de verdad esto tiene que repetirse indefinidamente, o sea que cada vez que el thread esta por lockear, justo en esa ventana entre terminar el recorrido sin locks y efectivamente ejecutar `pred.lock()` y `curr.lock()`, aparece un nuevo nodo insertado entre su `pred` y su `curr`. Si la insercion ocurriera despues de que el thread ya tomo ambos locks no habria problema, porque el otro thread que quiere insertar se bloquearia esperando esos mismos locks.

Un ejemplo concreto: lista `head -> A(1) -> C(5) -> null`. El thread T1 quiere hacer `remove` del item con `key=5` (osea eliminar a `C`). T1 recorre sin locks y determina `pred = A`, `curr = C`. Antes de que T1 llegue a `pred.lock()`, otro thread T2 inserta un nodo `B(3)` entre `A` y `C`, quedando la lista `A -> B -> C`. T1 sigue con sus referencias viejas (`pred=A`, `curr=C`), toma los locks de `A` y `C`, y llama a `validate(A, C)`, que falla porque ahora `A.next == B`, no `C`. T1 reintenta desde cero: recorre de nuevo, encuentra `pred=B`, `curr=C`, pero justo antes de lockear, un thread T3 inserta otro nodo `D(4)` entre `B` y `C`. Vuelve a fallar. Si este patron se repite indefinidamente, con distintos threads insertando justo en el momento critico cada vez, T1 nunca logra eliminar a `C`.

## Punto B 
Si se cambiara el orden en el que se agarran los locks y primero se agarra `curr` y luego `pred` en `add` pero en `remove` sigue el mismo orden de agarrar los locks de primero agarrar `pred` y despues `current`, puede haber un deadlock si por ejemplo tenemos los nodos `head -> A(1) -> B(5) -> D(10) -> E(15) -> null`

Si un thread `T1` quiere agregar un elemento `C(9)` que deberia ir luego de `D(10)`, y otro thread quiere eliminar `D(10)`, podria ocurrir que el add empieze ejecutandose y tome el lock de `D(10)` luego sea cortado antes de tomar el lock de `B(5)`, luego el otro thread como primero toma el lock del predecesor y luego el del sucesor, va a tomar el de `B(5)`, y cuando quiera tomar el de `D(10)` va a haber una espera circular que termina en deadlock. 

Lo importante es que en todos los metodos se tomen los locks en el mismo orden para que no hayan dependencias circulares que provoquen deadlocks. 

# Ejercicio 5 
## Punto A 
```java 
public class Tabla{
    HashMap<Integer, Object> tabla = new HashMap <Integer, Object>();
    private final Lock lock = new ReentrantLock();
    private void actualizarEntrada(int i){
        while(true){
            try{
                Object v1 = tabla.get(i);
                Object v2 = ComputacionalmenteCostoso(v1);
                lock.lock();
                if(tabla.get(i) == v1){
                    tabla.remove(i);
                    tabla.put(i, v2);
                    return;
                }
            }
            finally{
                lock.unlock();
            }
        }
    }
}
```
## Punto B 
La solucion propuesta no necesariamente es libre de inanicion, si un thread trata de modificar un valor y siempre le ganan  y modifican ese valor antes de que compruebe, va a volver a intentar indefinidamente hasta que pueda cambiar el valor, si siempre le ganan va a tener inanicion. 

# Ejercicio 6
## Punto A
Por mas que `in` y `out` esten dentro de monitores, al estar trackeando el estado de la cola partiendo el estado en 2 podemos tener problemas. 

Si primero un thread `T1` llama a `enq` van a quedar los stacks `in` con el elemento encolado `c1` y el stack `out` vacio. 

Si luego `T1` y `T2` tratan de desencolar concurrentemente ese unico elemento puede ocurrir que: 
- `T1` empieza a ejecutar `deq`, como la cola `out` esta vacia va a entrar al cuerpo del `if`. 
- `T2` hace lo mismo y entra tambien al if. 
- `T1` ejecuta una vez el ciclo vaciando `in` y poniendo un elemento en `out`. Sale del ciclo y antes de calcular `out.pop()` para devolverlo es context switcheado. 
- `T2` no ejecuta el ciclo porque ya se cumple que `in` esta vacio, asi que llega a `out.pop()` y devuelve. 
- `T1` al intentar ejecutar `out.pop()` ese stack ya esta vacio, asi que va a ocurrir una excepcion en tiempo de ejecucion. 

Otro problema que puede ocurrir es que no se cumpla la condicion de FIFO. Si hay una cola con un solo elemento `C1`, si concurrentemente un thread `T1` ejecuta `deq()` mientras otro thread `T2` ejecuta `enq(C2)` no se va a cumplir que va a sacar el primer elemento en entrar si no que va a sacar el segundo. 

Esto puede ocurrir si: 
- `T1` hace `enq(C1)` dejando `in` solo con el elemento `C1` y `out` vacio. 
- `T1` hace `deq()`, como `out` esta vacio, ejecuta el cuerpo del `if`, y como `in` no esta vacio va a pushear en out y sacar de in `C1`, antes de ver si tiene que ejecutar el ciclo es desalojado. 
- `T2` hace `enq(C2)`, haciendo que en `in` quede un stack que solo tiene a `C2`. 
- `T1` al revisar si `in` esta vacio ve que no y pushea en `out` `C2`. El problema es que ahora el orden del stack es el inverso al deseado en la cola, por lo que va a devolver `C2` en vez de `C1`. 

Tambien por como esta el codigo pueden ocurrir muchos problemas con la cola vacia al llamar a `deq` ya que no tiene ninguna excepcion en caso de que este vacia. 

## Punto B 
Asumiendo que nunca se va a llamar a `deq` con la cola vacia

```java 
class Stack {
    synchronized boolean isEmpty() { ... }
    synchronized Object pop() { ... }
    synchronized void push(Object o) { ... }
    ...
}

class Queue {
    Stack in  = new Stack();
    Stack out = new Stack();
    public void lock() { ... }
    public void unlock() { ... }

    void enq(Object o) {
        lock.lock();
        in.push(o);
        lock.unlock();
    }

    Object deq() {
        if (out.isEmpty()) {
            lock.lock();
            if(out.isEmpty()){ // Volvemos a chequear por si nos devuelven el lock tarde y ya pusieron elementos en out
                while (!in.isEmpty()) {
                    out.push(in.pop());
                }
            }
            
            lock.unlock();
        }
        return out.pop();
    }
}
```

# Ejercicio 7
La implementacion de una cola acotada que utiliza dos locks y la variable atomica `size` es: 
```java 
protected class Node {
    public T value;
    public volatile Node next;

    public Node(T x) {
        value = x;
        next = null;
    }
}

public class BoundedQueue<T> {
    ReentrantLock enqLock, deqLock;
    Condition notEmptyCondition;
    Condition notFullCondition;
    volatile int amountOfEnqs;
    volatile int amountOfDeqs;
    volatile Node head, tail;
    final int capacity;

    public BoundedQueue(int c) {
        capacity = c;
        head = new Node(null);
        tail = head;
        amountOfEnqs = 0;
        amountOfDeqs = 0;
        enqLock = new ReentrantLock();
        notFullCondition = enqLock.newCondition();
        deqLock = new ReentrantLock();
        notEmptyCondition = deqLock.newCondition();
    }

    public void enq(T x) {
        boolean mustWakeDequeuers = false;
        Node e = new Node(x);
        enqLock.lock();
        try {
            while (amountOfEnqs - amountOfDeqs == capacity) {
                notFullCondition.await();
            }
            tail.next = e;
            tail = e;
            amountOfEnqs += 1;
            if (amountOfEnqs - amountOfDeqs == 1) {
                mustWakeDequeuers = true;
            }
        } finally {
            enqLock.unlock();
        }
        if (mustWakeDequeuers) {
            deqLock.lock();
            try {
                notEmptyCondition.signalAll();
            } finally {
                deqLock.unlock();
            }
        }
    }

    public T deq() {
        T result;
        boolean mustWakeEnqueuers = false;
        deqLock.lock();
        try {
            while (head.next == null) {
                notEmptyCondition.await();
            }
            result = head.next.value;
            head = head.next;
            amountOfDeqs += 1;
            if (amountOfEnqs - amountOfDeqs == capacity - 1) {
                mustWakeEnqueuers = true;
            }
        } finally {
            deqLock.unlock();
        }
        if (mustWakeEnqueuers) {
            enqLock.lock();
            try {
                notFullCondition.signalAll();
            } finally {
                enqLock.unlock();
            }
        }
        return result;
    }
}
```

# Ejercicio 8 
Una primer implementacion puede ser: 
```java
protected class Node {
    public T value;
    public volatile Node next;

    public Node(T x) {
        value = x;
        next = null;
    }
}

public class BoundedStack<T> {
    ReentrantLock stackLock;
    Condition notEmptyCondition;
    Condition notFullCondition;
    int size;
    volatile Node top;
    final int capacity;

    public BoundedStack(int c) {
        capacity = c;
        top = null;
        size = 0;
        stackLock = new ReentrantLock();
        notFullCondition = stackLock.newCondition();
        notEmptyCondition = stackLock.newCondition();
    }

    public void push(T x) {
        Node e = new Node(x);
        stackLock.lock();
        try {
            while (size == capacity) {
                notFullCondition.await();
            }
            e.next = top;
            top = e;
            size++;
            notEmptyCondition.signalAll();
        } finally {
            stackLock.unlock();
        }
    }

    public T pop() {
        T result;
        stackLock.lock();
        try {
            while (top == null) {
                notEmptyCondition.await();
            }
            result = top.value;
            top = top.next;
            size--;
            notFullCondition.signalAll();
        } finally {
            stackLock.unlock();
        }
        return result;
    }
}
```

La implementacion libre de locks con backoffs va a ser: 
```java 
public class Node<T> {
    public T value;
    public Node<T> next;

    public Node(T value) {
        this.value = value;
        this.next = null;
    }
}

public class LockFreeStack<T> {
    AtomicReference<Node<T>> top = new AtomicReference<Node<T>>(null);

    static final int MIN_DELAY = ...;
    static final int MAX_DELAY = ...;
    final int capacity;
    Backoff backoff = new Backoff(MIN_DELAY, MAX_DELAY);
    AtomicInteger size = new AtomicInteger(0);

    public LockFreeStack(int capacity) {
        this.capacity = capacity;
    }

    protected boolean tryPush(Node<T> node) {
        if (size.get() == capacity) {
            return false;
        }
        Node<T> oldTop = top.get();
        node.next = oldTop;
        if (top.compareAndSet(oldTop, node)) {
            size.getAndIncrement();
            return true;
        }
        return false;
    }

    public void push(T value) {
        Node<T> node = new Node<T>(value);
        while (true) {
            if (tryPush(node)) {
                return;
            } else {
                backoff.backoff();
            }
        }
    }

    protected Node<T> tryPop() throws EmptyException {
        Node<T> oldTop = top.get();
        if (oldTop == null) {
            throw new EmptyException();
        }
        Node<T> newTop = oldTop.next;
        if (top.compareAndSet(oldTop, newTop)) {
            size.getAndDecrement();
            return oldTop;
        } else {
            return null;
        }
    }

    public T pop() throws EmptyException {
        while (true) {
            Node<T> returnNode = tryPop();
            if (returnNode != null) {
                return returnNode.value;
            } else {
                backoff.backoff();
            }
        }
    }
}
```

# Ejercicio 9 
TODO 

# Ejercicio 10 
Solucion `lock-free` usando una unica variable atomica: 
```java 
public class Contador{
    AtomicInteger value = new AtomicInteger(0);

    public void inc(){
        value.getAndIncrement();
    }

    public int get(){
        return value.get();
    }

    public void reset(){
        value.getAndSet(0);
    }
}
```
Esta solucion no es `wait-free` porque tiene un alto nivel de contencion en la variable atomica `value`, por mas que sea atomica puede ocurrir que si un proceso quiere hacer alguna interaccion con la variable, y constantemente lleguen otros threads a llamar a un metodo no lo seleccionen para ejecutar. Por lo que no podemos determinar en que cantidad finita de pasos un thread puede interactuar con `value`. 

# Ejercicio 11