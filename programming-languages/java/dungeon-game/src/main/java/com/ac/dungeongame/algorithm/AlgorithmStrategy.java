package com.ac.dungeongame.algorithm;

import com.ac.dungeongame.model.Game;

public interface AlgorithmStrategy {
    Game calculate(int[][] dungeon);
}
