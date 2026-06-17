package com.ac.chat.domain;

import java.time.LocalDateTime;

public record ChatMessage (Long id, String content, LocalDateTime createdAt, Author author){
    public ChatMessage(String content, Author author) {
        this(IdGenerator.getNextId("message"), content, LocalDateTime.now(), author);
    }
}
