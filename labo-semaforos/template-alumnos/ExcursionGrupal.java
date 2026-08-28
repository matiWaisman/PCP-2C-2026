// Volcan Lanin, parte 2.

public class MainExcursionGrupal {
    static int PORTENIOS = 6;
    static int STAGES = 4;

    // TODO: variables/objetos de sincronizacion.
    static int cantidad_llegados = 0; 
    private semaphore mutex = Semaphore(0);
    private semaphore s1 = Semaphore(0);
    private semaphore s2 = Semaphore(0);
    public void esperarPirca() throws InterruptedException {
        // TODO
        mutex.acquire();
        cantidad_llegados++; 
        if (cont == PORTENIOS){
            s1.release(PORTENIOS);
            cantidad_llegados = 0;
        }
        mutex.release();
        s1.acquire();
        m.acquire();
        cantidad_llegados++; 
        if(cantidad_llegados == PORTENIOS){
            s2.release(PORTENIOS);
            cantidad_llegados = 0;
        }


    }

    public static void main(String[] args) throws InterruptedException {
        MainExcursionGrupal excursion = new MainExcursionGrupal();

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
