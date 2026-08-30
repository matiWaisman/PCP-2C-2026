// Lollapallozers con acceso justo mediante un molinete FIFO.
import java.util.concurrent.Semaphore;
import java.util.concurrent.ThreadLocalRandom;

public class VentaTicketsFila {
    static int TYPES = 3;
    static int[] LIMIT = { 50, 30, 10 };
    static int BUYERS = 20;
    static int ATTEMPTS_PER_BUYER = 500;

    private int[] available = LIMIT.clone();
    private int cantidadLeyendo = 0;

    private final Semaphore mutexLectura = new Semaphore(1);
    private final Semaphore acceso = new Semaphore(1);
    private final Semaphore molinete = new Semaphore(1, true);

    // Devuelve true si logro comprar un ticket del tipo dado, false si no quedan.
    public boolean comprar(int type) throws InterruptedException {
        while (true) {
            // Entrada al grupo de lectores a traves del molinete FIFO.
            molinete.acquire();
            mutexLectura.acquire();
            cantidadLeyendo++;
            if (cantidadLeyendo == 1) {
                acceso.acquire();
            }
            mutexLectura.release();
            molinete.release();

            int cantidadLeida = available[type];

            // Salida del grupo de lectores.
            mutexLectura.acquire();
            cantidadLeyendo--;
            if (cantidadLeyendo == 0) {
                acceso.release();
            }
            mutexLectura.release();

            if (cantidadLeida == 0) {
                return false;
            }

            // El comprador conserva su lugar en la fila mientras espera
            // que terminen los lectores que ya habian ingresado.
            molinete.acquire();
            acceso.acquire();
            molinete.release();

            if (cantidadLeida == available[type]) {
                available[type]--;
                acceso.release();
                return true;
            }

            acceso.release();
        }
    }

    public static void main(String[] args) throws InterruptedException {
        VentaTicketsFila sale = new VentaTicketsFila();
        int[] soldPerBuyer = new int[BUYERS];

        Thread[] buyers = new Thread[BUYERS];
        for (int i = 0; i < BUYERS; i++) {
            int buyer = i;
            buyers[i] = new Thread(() -> {
                try {
                    for (int j = 0; j < ATTEMPTS_PER_BUYER; j++) {
                        int typeToBuy = ThreadLocalRandom.current().nextInt(TYPES);
                        if (sale.comprar(typeToBuy)) {
                            soldPerBuyer[buyer]++;
                        }
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
        }

        for (Thread t : buyers) t.start();
        for (Thread t : buyers) t.join();

        int totalSold = 0;
        for (int v : soldPerBuyer) totalSold += v;

        int totalRemaining = 0;
        boolean oversold = false;
        for (int type = 0; type < TYPES; type++) {
            int remaining = sale.available[type];
            totalRemaining += remaining;
            if (remaining < 0) oversold = true;
        }

        int totalLimit = 0;
        for (int limit : LIMIT) totalLimit += limit;

        System.out.println("Vendidos: " + totalSold + ", restante: " + totalRemaining + ", total: " + totalLimit);
        boolean ok = !oversold && (totalSold + totalRemaining == totalLimit);
        System.out.println(ok ? "OK" : "MAL");
    }
}
