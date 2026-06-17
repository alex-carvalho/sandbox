package com.ac.chat.repository;

import com.ac.chat.domain.ChatSession;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Component
public class ChatRepository {

    private final Map<Long, ChatSession> sessions = new ConcurrentHashMap<>();

    public Flux<ChatSession> findAll() {
        return Flux.fromIterable(sessions.values());
    }

    public Mono<ChatSession> save(ChatSession chatSession) {
        sessions.put(chatSession.id(), chatSession);
        return Mono.just(chatSession);
    }

    public Mono<ChatSession> findById(Mono<Long> id) {
        return id.map(sessions::get);
    }

}
