package com.example.dungeongame.controller;

import com.example.dungeongame.model.Game;
import com.example.dungeongame.service.GameService;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;

@SpringBootTest(properties = {"spring.profiles.active=test"})
@AutoConfigureMockMvc
class GameControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    void testCalculateMinimumHPWithABTest() throws Exception {
        String dungeonJson = "{\"dungeon\": [[-2, -3, 3], [-5, -10, 1], [10, 30, -5]]}";

        mockMvc.perform(post("/game")
                .contentType(MediaType.APPLICATION_JSON)
                .content(dungeonJson))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.minimumHp").isNumber());
    }
}
