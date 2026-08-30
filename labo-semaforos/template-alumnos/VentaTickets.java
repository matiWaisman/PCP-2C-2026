// Lollapallozers.
import java.util.concurrent.Semaphore;
import java.util.concurrent.ThreadLocalRandom;

public class VentaTickets {
    static int TYPES = 3;
    static int[] LIMIT = { 50, 30, 10 };
    static int BUYERS = 20;
    static int ATTEMPTS_PER_BUYER = 500;

    private int[] available = LIMIT.clone();
    private int cantidadEsperandoALeer = 0;
    private int cantidadEsperandoAComprar = 0;
    private Semaphore mutexLectura = new Semaphore(1);
    private Semaphore mutexCompras = new Semaphore(1);
    private Semaphore puertaLectores = new Semaphore(1);
    private Semaphore acceso = new Semaphore(1);


    // Solucion con prioridad compradores

    // Devuelve true si logro comprar un ticket del tipo dado, false si no quedan.
    public boolean comprar(int type) throws InterruptedException {
        while(true){
            puertaLectores.acquire();
            mutexLectura.acquire();
            cantidadEsperandoALeer += 1;
            if(cantidadEsperandoALeer == 1){
                acceso.acquire();
            }
            mutexLectura.release();
            puertaLectores.release();
            // Si estamos aca es pq nadie esta comprando todavia, si hubiera alguien comprando estariamos colgados en acceso.aqcuire()
            int cantidad = available[type];
            // Ya leimos, asi que nos bajamos de la cantidad esperando a leer 
            mutexLectura.acquire();
            cantidadEsperandoALeer -= 1;
            if(cantidadEsperandoALeer == 0){
                acceso.release(); // El que quiera comprar o leer va a poder
            }
            mutexLectura.release();
            if(cantidad == 0){
                return false;
            }
            // Ahora queremos comprar asi que pedimos acceso exclusivo
            // Primero vemos cuantos estan esperando para comprar
            mutexCompras.acquire();
            cantidadEsperandoAComprar += 1; 
            if(cantidadEsperandoAComprar == 1){
                // Cerramos la puerta de lectura consumiendo un slot 
                puertaLectores.acquire();
            }
            mutexCompras.release();
            acceso.acquire();
            if(cantidad == available[type]){
                available[type] -= 1;
                mutexCompras.acquire();
                cantidadEsperandoAComprar -= 1;
                if(cantidadEsperandoAComprar == 0){
                    puertaLectores.release();
                }
                mutexCompras.release();
                acceso.release();
                return true;
            }
            // Si no terminamos tenemos que volver a intentar
            mutexCompras.acquire();
            cantidadEsperandoAComprar -= 1;
            if(cantidadEsperandoAComprar == 0){
                puertaLectores.release();
            }
            mutexCompras.release();
            acceso.release();
        }
    }

    public static void main(String[] args) throws InterruptedException {
        VentaTickets sale = new VentaTickets();
        int[] soldPerBuyer = new int[BUYERS];

        // Inicializacion de los Threads 
        Thread[] buyers = new Thread[BUYERS];
        for (int i = 0; i < BUYERS; i++){
            int buyer = i;
            buyers[i] = new Thread(()-> {
                try{
                    for (int j = 0; j < ATTEMPTS_PER_BUYER; j++){
                        int type_to_buy = ThreadLocalRandom.current().nextInt(TYPES);
                        if(sale.comprar(type_to_buy)){
                            soldPerBuyer[buyer]++;
                        }
                    }
                } catch (InterruptedException e){
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
        for (int l : LIMIT) totalLimit += l;

        System.out.println("Vendidos: " + totalSold + ", restante: " + totalRemaining + ", total: " + totalLimit);
        boolean ok = !oversold && (totalSold + totalRemaining == totalLimit);
        System.out.println(ok ? "OK" : "MAL");
    }
}
