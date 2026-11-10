import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.RecursiveAction;

public class TransposeTask extends RecursiveAction {

    private final int[][] matrix;
    private final int rowStart;
    private final int rowEnd;
    private final int colStart;
    private final int colEnd;

    // Filas: [rowStart, rowEnd), columnas: [colStart, colEnd).
    // Los inicios están incluidos y los finales están excluidos.
    public TransposeTask(int[][] matrix, int rowStart, int rowEnd,
                         int colStart, int colEnd) {
        this.matrix = matrix;
        this.rowStart = rowStart;
        this.rowEnd = rowEnd;
        this.colStart = colStart;
        this.colEnd = colEnd;
    }

    @Override
    public void compute() {
        int cantidadElementos = (rowEnd - rowStart) * (colEnd - colStart);
        if(cantidadElementos == 4){
            for(int i = rowStart; i < rowEnd; i++){
                for(int j = colStart; j < colEnd; j++){
                    if(i < j){
                        int aux = matrix[i][j];
                        matrix[i][j] = matrix[j][i];
                        matrix[j][i] = aux;
                    }
                }
            }
            return;
        }
        List <TransposeTask> tasks = new ArrayList<>(4);
        int mitadFilas = rowStart + (rowEnd - rowStart) / 2;
        int mitadColumnas = colStart + (colEnd - colStart) / 2;
        
        tasks.add(new TransposeTask(matrix, rowStart, mitadFilas, colStart, mitadColumnas));
        tasks.add(new TransposeTask(matrix, rowStart, mitadFilas, mitadColumnas, colEnd));
        tasks.add(new TransposeTask(matrix, mitadFilas, rowEnd, colStart, mitadColumnas));
        tasks.add(new TransposeTask(matrix, mitadFilas, rowEnd, mitadColumnas, colEnd));
        
        tasks.stream().forEach((task)->{task.fork();});
        tasks.stream().forEach((task)->{task.join();});   
    }

}
