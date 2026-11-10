import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.RecursiveTask;

public class VectorSumTask extends RecursiveTask<Integer> {

    private final int[] vector;
    private final int left;
    private final int right;

    // Rango: [left, right)
    // left incluido, right excluido.
    public VectorSumTask(int[] vector, int left, int right) {
        this.vector = vector;
        this.left = left;
        this.right = right;
    }

    @Override
    protected Integer compute() {
        int cantidad = right - left;
        if (cantidad == 2) {
            return vector[left] + vector[left + 1];
        }
        int middle = left + (right - left) / 2;
        VectorSumTask rightTask = new VectorSumTask(vector, middle, right);
        rightTask.fork();
        VectorSumTask leftTask = new VectorSumTask(vector, left, middle);

        return leftTask.compute() + rightTask.join();
    }

    public static void main(String[] args) {
        // Longitudes potencia de 2 y mayores o iguales a 2.
        int[][] casos = {
            {1, 2},
            {1, 2, 3, 4},
            {-3, 5, -2, 4, 0, 6, -1, 7},
            new int[8],
            new int[1024]
        };

        for (int i = 0; i < casos[4].length; i++) {
            casos[4][i] = i % 21 - 10;
        }

        ForkJoinPool pool = new ForkJoinPool(2);
        try {
            for (int i = 0; i < casos.length; i++) {
                int[] vector = casos[i];

                int esperado = 0;
                for (int valor : vector) {
                    esperado += valor;
                }

                int obtenido = pool.invoke(
                    new VectorSumTask(vector, 0, vector.length)
                );

                System.out.printf(
                    "Caso %d (n=%d): esperado=%d, obtenido=%d%n",
                    i + 1, vector.length, esperado, obtenido
                );

                if (obtenido != esperado) {
                    throw new AssertionError("Fallo en el caso " + (i + 1));
                }
            }

            System.out.println("Todos los casos pasaron.");
        } finally {
            pool.shutdown();
        }
    }
}
