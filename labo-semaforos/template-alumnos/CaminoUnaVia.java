// Zorros grises de Saldungaray, parte 1 (ghostbusters).
import java.util.concurrent.Semaphore;
import java.util.concurrent.ThreadLocalRandom;

// Idea: Hacer que se vayan alternando de a K autos de cada lado. 
// Si toca el turno de alternar pero del otro lado no hay autos, seguis vos con K autos.

public class CaminoUnaVia {
    static int NUM_CARS = 12;
    
    private final int[] cantidadEsperando;
    private final Semaphore[] s;
    private final Semaphore[] mutexs;
    private Semaphore turnstile = new Semaphore(1, true);
    private Semaphore cambioSentido = new Semaphore(1);

    public CaminoUnaVia() {
        s = new Semaphore[2];
        mutexs = new Semaphore[2];
        cantidadEsperando = new int[2];

        for (int direction = 0; direction < 2; direction++) {
            s[direction] = new Semaphore(0);
            mutexs[direction] = new Semaphore(1);
            cantidadEsperando[direction] = 0;
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
    }

    public void salir(int direction) throws InterruptedException {
        mutexs[direction].acquire();
        cantidadEsperando[direction] -= 1; 
        if (cantidadEsperando[direction] == 0){
            cambioSentido.release();
        } 
        mutexs[direction].release();
    }

    // Simula el tiempo que tarda un auto en cruzar el desvio.
    private void cruzarPuente() throws InterruptedException {
        Thread.sleep(ThreadLocalRandom.current().nextInt(2));
    }

    public static void main(String[] args) throws InterruptedException {
        CaminoUnaVia route = new CaminoUnaVia();
        Thread[] autos = new Thread[NUM_CARS];
        for (int i = 0; i < NUM_CARS; i++) {
            int car = i;
            int direction = ThreadLocalRandom.current().nextInt(2);
            autos[i] = new Thread(() -> {
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
