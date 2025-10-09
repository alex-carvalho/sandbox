import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import java.util.List;

class BitonicSequenceTest {

    @Test
    void testGenerateLargeSequence() {
        List<Integer> sequence = BitonicSequence.generate(5,3,10);
        List<Integer> sequence2 = BitonicSequence.generate(7, 2, 5);

        assertAll( () -> assertEquals(List.of(9, 10, 9, 8, 7), sequence),
                () -> assertEquals(List.of(2, 3, 4, 5, 4, 3, 2), sequence2) );
    }
}