package com.ac.chat.domain;

import java.time.LocalDateTime;
import java.util.List;

public record ChatSession(Long id, String title, LocalDateTime createdAt, List<ChatMessage> messages){
    public ChatSession(String title, List<ChatMessage> messages) {
        this(IdGenerator.getNextId("session"), title, LocalDateTime.now(), messages);
    }
}