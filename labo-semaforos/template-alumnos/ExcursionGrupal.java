// Volcan Lanin, parte 2.
import java.util.concurrent.Semaphore;

public class ExcursionGrupal {
    static int PORTENIOS = 6;
    static int STAGES = 4;

    static int cantidad_llegados = 0; 
    private Semaphore mutex = new Semaphore(1);
    private Semaphore s1 = new Semaphore(0);
    private Semaphore s2 = new Semaphore(0);
    public void esperarPirca() throws InterruptedException {
        mutex.acquire();
        cantidad_llegados++; 
        if (cantidad_llegados == PORTENIOS){
            s1.release(PORTENIOS);
            cantidad_llegados = 0;
        }
        mutex.release();
        s1.acquire();
        mutex.acquire();
        cantidad_llegados++; 
        if(cantidad_llegados == PORTENIOS){
            s2.release(PORTENIOS);
            cantidad_llegados = 0;
        }
        mutex.release();
        s2.acquire();
    }

    public static void main(String[] args) throws InterruptedException {
        ExcursionGrupal excursion = new ExcursionGrupal();

        Thread[] group = new Thread[PORTENIOS];
        for (int i = 0; i < PORTENIOS; i++) {
            int id = i;
            group[i] = new Thread(() -> {
                try {
                    for (int stage = 1; stage <= STAGES; stage++) {
                        Thread.sleep((long) (Math.random() * 300));
                        System.out.println("porteno " + id + " llega a la pirca del tramo " + stage);
                        excursion.esperarPirca();
                        System.out.println("porteno " + id + " arranca el tramo " + (stage + 1));
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
        }

        for (Thread t : group) t.start();
        for (Thread t : group) t.join();
    }
}
