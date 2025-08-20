package com.example.dungeongame.service;

import com.example.dungeongame.model.Game;
import com.example.dungeongame.repository.GameRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.Optional;

@Service
public class GameService {

    @Autowired
    private GameRepository gameRepository;

    private Game calculateMinimumHP(int[][] dungeon) {
        Game game = new Game();
        game.setDungeon(Arrays.deepToString(dungeon));

        int m = dungeon.length;
        int n = dungeon[0].length;

        // dp[i][j] represents the minimum health required to enter room (i, j)
        int[][] dp = new int[m][n];

        // Calculate health required for the princess's room
        dp[m - 1][n - 1] = Math.max(1, 1 - dungeon[m - 1][n - 1]);

        // Calculate health required for the last row (moving from right to left)
        for (int i = n - 2; i >= 0; i--) {
            dp[m - 1][i] = Math.max(1, dp[m - 1][i + 1] - dungeon[m - 1][i]);
        }

        // Calculate health required for the last column (moving from bottom to top)
        for (int i = m - 2; i >= 0; i--) {
            dp[i][n - 1] = Math.max(1, dp[i + 1][n - 1] - dungeon[i][n - 1]);
        }

        // Calculate health required for the rest of the rooms
        for (int i = m - 2; i >= 0; i--) {
            for (int j = n - 2; j >= 0; j--) {
                int minHealthOnExit = Math.min(dp[i + 1][j], dp[i][j + 1]);
                dp[i][j] = Math.max(1, minHealthOnExit - dungeon[i][j]);
            }
        }

        game.setMinimumHp(dp[0][0]);
        return gameRepository.save(game);
    }

    private Game calculateMinimumHPOptimized(int[][] dungeon) {
        Game game = new Game();
        game.setDungeon(Arrays.deepToString(dungeon));

        int m = dungeon.length;
        int n = dungeon[0].length;

        // Optimized algorithm using a single-dimensional array
        int[] dp = new int[n];
        dp[n - 1] = Math.max(1, 1 - dungeon[m - 1][n - 1]);

        for (int i = n - 2; i >= 0; i--) {
            dp[i] = Math.max(1, dp[i + 1] - dungeon[m - 1][i]);
        }

        for (int i = m - 2; i >= 0; i--) {
            dp[n - 1] = Math.max(1, dp[n - 1] - dungeon[i][n - 1]);
            for (int j = n - 2; j >= 0; j--) {
                dp[j] = Math.max(1, Math.min(dp[j], dp[j + 1]) - dungeon[i][j]);
            }
        }

        game.setMinimumHp(dp[0]);
        return gameRepository.save(game);
    }

    public Game calculate(int[][] dungeon) {
        boolean useOptimized = Math.random() < 0.5; // 50% chance for A or B
        if (useOptimized) {
            return calculateMinimumHPOptimized(dungeon);
        } else {
            return calculateMinimumHP(dungeon);
        }
    }

    public Game saveGame(Game game) {
        return gameRepository.save(game);
    }

    public Optional<Game> findGameById(Long id) {
        return gameRepository.findById(id);
    }
}