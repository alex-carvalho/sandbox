package com.ac.dungeongame.algorithm;

import com.ac.dungeongame.model.Game;
import org.springframework.stereotype.Component;
import java.util.Arrays;

@Component("Optimized")
public class OptimizedAlgorithm implements AlgorithmStrategy {

    @Override
    public Game calculate(int[][] dungeon) {
        Game game = new Game();
        game.setDungeon(Arrays.deepToString(dungeon));

        int m = dungeon.length;
        int n = dungeon[0].length;

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
        return game;
    }
}
