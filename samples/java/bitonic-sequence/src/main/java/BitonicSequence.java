import java.util.*;

public class BitonicSequence {

    static List<Integer> generate(Integer n, Integer min, Integer max) {

        // If not possible
        if (n > (max - min) * 2 + 1) {
            return List.of(-1);
        }

        Deque<Integer> dq = new ArrayDeque<>();
        dq.addLast(max - 1);

        // Add decreasing part
        for (int i = max; i >= min && dq.size() < n; i--) {
            dq.addLast(i);
        }

        // Add increasing part
        for (int i = max - 2; i >= min && dq.size() < n; i--) {
            dq.addFirst(i);
        }

        return new ArrayList<>(dq);
    }
    

    public static void main(String[] args) {
        // Generate bitonic sequence of length 8 from range [1, 20]
        List<Integer> sequence = generate(8, 1, 20);
        System.out.println("Bitonic sequence: " + sequence);
        
        // Generate another example
        sequence = generate(10, 5, 50);
        System.out.println("Another bitonic sequence: " + sequence);
    }
}