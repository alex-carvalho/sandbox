package com.ac.chat.controller;

import com.ac.chat.domain.ChatRequest;
import com.ac.chat.domain.ChatSession;
import com.ac.chat.service.ChatService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@RestController
public class ChatController {
    private static final Logger log = LoggerFactory.getLogger(ChatController.class);

    private final ChatService chatService;

    public ChatController(ChatService chatService) {
        this.chatService = chatService;
    }

    @GetMapping("/sessions")
    public Flux<ChatSession> getSessions() {
        log.info("Received GET /sessions");
        return chatService.getSessions();
    }

    @GetMapping("/sessions/{id}")
    public Mono<ChatSession> getSession(@PathVariable Long id) {
        log.info("Received GET /sessions/{}", id);
        return chatService.getSession(Mono.justOrEmpty(id));
    }

    @PostMapping(value = "/chat", consumes =  MediaType.APPLICATION_JSON_VALUE)
    public Mono<String> chat(@RequestBody ChatRequest request) {
        log.info("Received POST /chat conversationId={} prompt=\"{}\"",
                request.conversationId(), request.prompt());
        return chatService.chat(request);
    }

    @PostMapping(value = "/chat/stream", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<String> streamChat(@RequestBody ChatRequest request) {
        log.info("Received POST /chat/stream conversationId={}  prompt=\"{}\"",
                request.conversationId(), request.prompt());
        return chatService.streamChat(request)
                .onErrorResume(e -> {
                    log.error("POST /chat/stream failed conversationId={}", request.conversationId(), e);
                    return Flux.just("Error: " + e.getMessage());
                });
    }

}
