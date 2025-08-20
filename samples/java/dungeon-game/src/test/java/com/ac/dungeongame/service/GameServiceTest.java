package com.ac.dungeongame.service;

import com.ac.dungeongame.model.Game;
import com.ac.dungeongame.repository.GameRepository;
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
        Game savedGame = gameService.saveGame(game);

        assertNotNull(savedGame.getId());
    }

    @Test
    void testFindGameById() {
        Game game = new Game();
        Game savedGame = gameRepository.save(game);

        Optional<Game> foundGame = gameService.findGameById(savedGame.getId());
        assertTrue(foundGame.isPresent());
    }

    @Test
    void testCalculateMinimumHP() {
        int[][] dungeon = {
            {-2, -3, 3},
            {-5, -10, 1},
            {10, 30, -5}
        };

        Game result = gameService.calculate(dungeon);

        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals("[[-2, -3, 3], [-5, -10, 1], [10, 30, -5]]", result.getDungeon());
        assertTrue(result.getMinimumHp() > 0, "Minimum HP should be greater than 0");
    }

    @Test
    void testCalculateMinimumHPWithDifferentDungeon() {
        int[][] dungeon = {
            {0, 0, 0},
            {0, 0, 0},
            {0, 0, 0}
        };

        Game result = gameService.calculate(dungeon);

        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals("[[0, 0, 0], [0, 0, 0], [0, 0, 0]]", result.getDungeon());
        assertEquals(1, result.getMinimumHp(), "Minimum HP should be 1 for an empty dungeon");
    }
}