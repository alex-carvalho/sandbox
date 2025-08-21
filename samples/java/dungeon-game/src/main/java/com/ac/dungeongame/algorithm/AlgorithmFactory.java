package com.ac.dungeongame.algorithm;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Component;

@Component
public class AlgorithmFactory {

    @Autowired
    @Qualifier("Standard")
    private AlgorithmStrategy standardAlgorithm;

    @Autowired
    @Qualifier("Optimized")
    private AlgorithmStrategy optimizedAlgorithm;

    public AlgorithmStrategy getAlgorithm(AlgorithmType type) {
        switch (type) {
            case STANDARD:
                return standardAlgorithm;
            case OPTIMIZED:
                return optimizedAlgorithm;
            default:
                throw new IllegalArgumentException("Unknown algorithm type: " + type);
        }
    }
}
