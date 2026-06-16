package com.ac.simple;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

@RestController
public class JokeController {
    private final ChatClient chatClient;

    public JokeController(ChatClient.Builder chatClientBuilder) {
        this.chatClient = chatClientBuilder.build();
    }

    @GetMapping(value = "/joke", produces = MediaType.TEXT_PLAIN_VALUE)
    public Flux<String> chat() {
        return chatClient
                .prompt()
                .user("Tell me a joke")
                .stream()
                .content()
                .map(list -> String.join("", list));
    }
}
