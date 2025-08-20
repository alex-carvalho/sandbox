package com.ac.dungeongame.controller;

import com.ac.dungeongame.model.Game;
import com.ac.dungeongame.service.GameService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class GameController {

    @Autowired
    private GameService gameService;

    @PostMapping("/game")
    public Game calculateMinimumHP(@RequestBody DungeonRequest request) {
        return gameService.calculate(request.getDungeon());
    }
}

class DungeonRequest {
    private int[][] dungeon;

    public int[][] getDungeon() {
        return dungeon;
    }

    public void setDungeon(int[][] dungeon) {
        this.dungeon = dungeon;
    }
}
