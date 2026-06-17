package com.ac.chat.controller;

import com.ac.chat.domain.ChatRequest;
import com.ac.chat.domain.ChatSession;
import com.ac.chat.service.ChatService;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
public class ChatController {
    private final ChatService chatService;

    public ChatController(ChatService chatService) {
        this.chatService = chatService;
    }

    @GetMapping("/sessions")
    public Flux<ChatSession> getSessions() {
        return chatService.getSessions();
    }

    @GetMapping("/sessions/{id}")
    public Mono<ChatSession> getSession(@PathVariable Mono<Long> id) {
        return chatService.getSession(id);
    }

    @PostMapping(value = "/chat", consumes =  MediaType.APPLICATION_JSON_VALUE)
    public Mono<String> chat(@RequestBody ChatRequest request) {
        return chatService.chat(request);
    }

    @PostMapping(value = "/chat/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<String> streamChat(@RequestBody ChatRequest request) {
        return chatService.streamChat(request)
                .onErrorResume(e -> Flux.just("Error: " + e.getMessage()));
    }
}
