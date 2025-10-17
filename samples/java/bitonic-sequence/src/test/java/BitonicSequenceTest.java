import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import java.util.List;

class BitonicSequenceTest {

    @Test
    void testGenerateLargeSequence() {
        List<Integer> sequence = BitonicSequence.generate(5,3,10);
        List<Integer> sequence2 = BitonicSequence.generate(7, 2, 5);
        List<Integer> sequence3 = BitonicSequence.generate(6, 2, 5);
        List<Integer> sequence4 = BitonicSequence.generate(3, 1, 10);

        assertAll( () -> assertEquals(List.of(9, 10, 9, 8, 7), sequence),
                () -> assertEquals(List.of(2, 3, 4, 5, 4, 3, 2), sequence2),
                () -> assertEquals(List.of(3, 4, 5, 4, 3, 2), sequence3),
                () -> assertEquals(List.of(9,10,9), sequence4) );
    }
}