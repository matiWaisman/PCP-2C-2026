import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.RecursiveTask;

public class PowerTask extends RecursiveTask<Integer> {

    private final int x; 
    private final int c;

    public PowerTask(int x, int c) {
        this.x = x;
        this.c = c;
    }

    @Override
    protected Integer compute() {
        if(c == 0){
            return 1;
        }
        if(c == 1){
            return x;
        }
        if(c == 2){
            return x * x;
        }
        PowerTask task = new PowerTask(x, c/2);
        task.fork();
        int res = task.join();
        if(c % 2 == 1){
            return res * res * x;
        }
        return res * res;
    }

}
