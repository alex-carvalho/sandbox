package com.example.dungeongame.controller;

import com.example.dungeongame.model.Game;
import com.example.dungeongame.service.GameService;
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
        return gameService.calculateMinimumHP(request.getDungeon());
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
