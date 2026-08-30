import java.util.concurrent.Semaphore;

public class ComputarNImpares{
    int n = 1;
    int i = 1;
    int limite;
    int suma = 0;

    private Semaphore puedeGenerar = new Semaphore(0);
    private Semaphore puedeConsumir = new Semaphore(1);

    public void setLimite(int N){
        this.limite = N;
    }

    public void generador() throws InterruptedException{
        while(true){
            puedeGenerar.acquire();
            if(i <= limite){
                n = n + 2;
                i = i + 1;
                puedeConsumir.release();
            }
            else{
                System.out.println(suma);
                puedeConsumir.release();
                break;
            }
        }
    }

    public void consumidor() throws InterruptedException{
        while(true){
            puedeConsumir.acquire();
            if(i > limite){
                puedeGenerar.release();
                break;
            }
            suma = suma + n;
            puedeGenerar.release();
        }
    }

    public static void main(String[] args) throws InterruptedException{
        int N = Integer.parseInt(args[0]);
        ComputarNImpares computador = new ComputarNImpares();
        computador.setLimite(N);
        Thread generador = new Thread(() -> {
            try { 
                computador.generador();
            } 
            catch (InterruptedException e){
                Thread.currentThread().interrupt();
            }
        });

        Thread consumidor = new Thread(() -> {
            try { 
                computador.consumidor();
            } 
            catch (InterruptedException e){
                Thread.currentThread().interrupt();
            }
        });
        generador.start();
        consumidor.start();
        generador.join();
        consumidor.join();

    }
}