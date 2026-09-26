// Zorros grises de Saldungaray, parte 2 .
import java.util.concurrent.Semaphore;
import java.util.concurrent.ThreadLocalRandom;

public class CaminoUnaViaFIFO {
    static int NUM_CARS = 12;
    static int BRIDGE_CAPACITY = 4;

    // Cada auto es un thread y conserva su propio slot entre entrar() y salir().
    static class Auto extends Thread {
        private int miSlot;

        Auto(Runnable tarea) {
            super(tarea);
        }

        static Auto actual() {
            return (Auto) Thread.currentThread();
        }
    }

    private final int[] cantidadEsperando;
    private final Semaphore[] mutexs;
    private Semaphore turnstile = new Semaphore(1, true);
    private Semaphore cambioSentido = new Semaphore(1);


    private Semaphore cantidadMaximaAutos = new Semaphore(BRIDGE_CAPACITY);
    private Semaphore queueMutex = new Semaphore(1);
    private Semaphore[] exitQueue = new Semaphore[BRIDGE_CAPACITY];
    private int head = 0;
    private int nextSpot = 0;
    private int queued = 0;


    public CaminoUnaViaFIFO() {
        mutexs = new Semaphore[2];
        cantidadEsperando = new int[2];

        for (int direction = 0; direction < 2; direction++) {
            mutexs[direction] = new Semaphore(1);
            cantidadEsperando[direction] = 0;
        }

        for (int i = 0; i < BRIDGE_CAPACITY; i++) {
            exitQueue[i] = new Semaphore(0);
        }
    }

    public void entrar(int direction) throws InterruptedException {
        turnstile.acquire();
        mutexs[direction].acquire();
        cantidadEsperando[direction] += 1;
        if (cantidadEsperando[direction] == 1){
            cambioSentido.acquire();
        }
        mutexs[direction].release();
        turnstile.release();

        cantidadMaximaAutos.acquire();

        queueMutex.acquire();
        int slot = nextSpot;
        nextSpot = (slot + 1) % BRIDGE_CAPACITY;
        boolean esElPrimero = queued == 0;
        queued += 1;
        queueMutex.release();

        if(esElPrimero){
            exitQueue[slot].release();
        }

        Auto.actual().miSlot = slot;
    }

    public void salir(int direction) throws InterruptedException {
        exitQueue[Auto.actual().miSlot].acquire();

        mutexs[direction].acquire();
        cantidadEsperando[direction] -= 1; 
        if (cantidadEsperando[direction] == 0){
            cambioSentido.release();
        } 
        mutexs[direction].release();

        queueMutex.acquire();
        head = (head + 1) % BRIDGE_CAPACITY;
        queued--;
        boolean alguienQueda = queued > 0;
        int nuevoHead = head;
        queueMutex.release();

        if (alguienQueda) {
            exitQueue[nuevoHead].release();
        }

        cantidadMaximaAutos.release();

    }

    // Simula el tiempo que tarda un auto en cruzar el desvio.
    private void cruzarPuente() throws InterruptedException {
        Thread.sleep(ThreadLocalRandom.current().nextInt(2));
    }

    public static void main(String[] args) throws InterruptedException {
        CaminoUnaViaFIFO route = new CaminoUnaViaFIFO();
        Auto[] autos = new Auto[NUM_CARS];
        for (int i = 0; i < NUM_CARS; i++) {
            int car = i;
            int direction = ThreadLocalRandom.current().nextInt(2);
            autos[i] = new Auto(() -> {
                try{
                    route.entrar(direction);
                    System.out.println("Auto " + car + " cruzando en sentido " + direction);
                    route.cruzarPuente();
                    route.salir(direction);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }

            });

        }
        for (Thread t : autos) t.start();
        for (Thread t : autos) t.join();
    }
}
