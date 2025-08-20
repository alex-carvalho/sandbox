package com.example.dungeongame.service;

import com.example.dungeongame.model.Game;
import com.example.dungeongame.repository.GameRepository;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@ActiveProfiles("test")
class GameServiceTest {

    @Autowired
    private GameService gameService;

    @Autowired
    private GameRepository gameRepository;

    @Test
    void testSaveGame() {
        Game game = new Game();
        game.setName("Test Game");
        Game savedGame = gameService.saveGame(game);

        assertNotNull(savedGame.getId());
        assertEquals("Test Game", savedGame.getName());
    }

    @Test
    void testFindGameById() {
        Game game = new Game();
        game.setName("Test Game");
        Game savedGame = gameRepository.save(game);

        Optional<Game> foundGame = gameService.findGameById(savedGame.getId());
        assertTrue(foundGame.isPresent());
        assertEquals("Test Game", foundGame.get().getName());
    }
}
