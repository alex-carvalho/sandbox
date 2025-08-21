package com.ac.dungeongame.service;

import com.ac.dungeongame.algorithm.AlgorithmFactory;
import com.ac.dungeongame.algorithm.AlgorithmStrategy;
import com.ac.dungeongame.algorithm.AlgorithmType;
import com.ac.dungeongame.model.AlgorithmExecution;
import com.ac.dungeongame.model.Game;
import com.ac.dungeongame.repository.AlgorithmExecutionRepository;
import com.ac.dungeongame.repository.GameRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import java.time.Duration;
import java.time.Instant;
import java.util.Optional;

@Service
public class GameService {

    @Autowired
    private GameRepository gameRepository;

    @Autowired
    private AlgorithmExecutionRepository algorithmExecutionRepository;

    @Autowired
    private AlgorithmFactory algorithmFactory;

    public Game calculate(int[][] dungeon) {
        AlgorithmType algorithmType = Math.random() < 0.5 ? AlgorithmType.OPTIMIZED : AlgorithmType.STANDARD;
        AlgorithmStrategy selectedAlgorithm = algorithmFactory.getAlgorithm(algorithmType);

        Instant start = Instant.now();
        Game game = selectedAlgorithm.calculate(dungeon);
        Instant end = Instant.now();
        long durationInNanos = Duration.between(start, end).toNanos();

        game = gameRepository.save(game); // Ensure the Game entity is saved before associating it

        AlgorithmExecution execution = new AlgorithmExecution();
        execution.setGame(game);
        execution.setAlgorithmType(algorithmType.name());
        execution.setDurationInNanos(durationInNanos);
        algorithmExecutionRepository.save(execution);

        return game;
    }

    public Game saveGame(Game game) {
        return gameRepository.save(game);
    }

    public Optional<Game> findGameById(Long id) {
        return gameRepository.findById(id);
    }
}