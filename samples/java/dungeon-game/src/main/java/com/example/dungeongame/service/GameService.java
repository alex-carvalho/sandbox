package com.example.dungeongame.service;

import com.example.dungeongame.model.Game;
import com.example.dungeongame.repository.GameRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Arrays;

@Service
public class GameService {

    @Autowired
    private GameRepository gameRepository;

    public Game calculateMinimumHP(int[][] dungeon) {
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
}