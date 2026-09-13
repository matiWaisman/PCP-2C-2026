import java.util.concurrent.Semaphore;
import java.util.Random;

public class ej7 {
    private static final int CANTIDAD_CLIENTES = 5;
    private static final int EJERCICIOS_POR_CLIENTE = 10;
    private static final int CANTIDAD_MAQUINAS = 4;

    private Semaphore[] mutexMaquinas = new Semaphore[CANTIDAD_MAQUINAS];
    private Semaphore rackDiscos;

    public ej7(int cantidadDiscos) {
        rackDiscos = new Semaphore(cantidadDiscos);
        for (int i = 0; i < CANTIDAD_MAQUINAS; i++) {
            mutexMaquinas[i] = new Semaphore(1); // 1 permiso = mutex binario
        }
    }

    public void usarMaquina(String nombreCliente, int indiceMaquina, int cantidadDiscos) throws InterruptedException{
        // Requiere que cantidadDiscos < cantidadDiscos
        mutexMaquinas[indiceMaquina].acquire();
        rackDiscos.acquire(cantidadDiscos);
        
        System.out.println("Usuario " + nombreCliente + " usa maquina " + indiceMaquina + " con cantidad de discos " + cantidadDiscos);
        Thread.sleep(1000);

        rackDiscos.release(cantidadDiscos);
        mutexMaquinas[indiceMaquina].release();
    }

    public static void main(String[] args) throws InterruptedException {
        int cantidadDiscosTotales = Integer.parseInt(args[0]);
        ej7 gimnasio = new ej7(cantidadDiscosTotales);

        Thread[] clientes = new Thread[CANTIDAD_CLIENTES];

        for (int i = 0; i < CANTIDAD_CLIENTES; i++) {
            String nombreCliente = "Cliente" + i;
            clientes[i] = new Thread(() -> {
                Random random = new Random();
                try {
                    for (int j = 0; j < EJERCICIOS_POR_CLIENTE; j++) {
                        int indiceMaquina = random.nextInt(CANTIDAD_MAQUINAS); // 0 a 3
                        int cantidadDiscos = random.nextInt(cantidadDiscosTotales) + 1; // 1 a cantidadDiscosTotales
                        gimnasio.usarMaquina(nombreCliente, indiceMaquina, cantidadDiscos);
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
        }

        for (Thread cliente : clientes) {
            cliente.start();
        }
        for (Thread cliente : clientes) {
            cliente.join();
        }

        System.out.println("Todos los clientes terminaron su rutina.");
    }

}