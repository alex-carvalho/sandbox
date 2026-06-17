package com.ac.chat.service;

import com.ac.chat.domain.Author;
import com.ac.chat.domain.ChatMessage;
import com.ac.chat.domain.ChatRequest;
import com.ac.chat.domain.ChatSession;
import com.ac.chat.repository.ChatRepository;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.messages.AssistantMessage;
import org.springframework.ai.chat.messages.Message;
import org.springframework.ai.chat.messages.UserMessage;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.List;

@Service
public class ChatService {
    private final ChatClient chatClient;
    private final ChatRepository chatRepository;

    public ChatService(ChatClient.Builder chatClientBuilder,
                       ChatRepository chatRepository) {
        this.chatClient = chatClientBuilder.build();
        this.chatRepository = chatRepository;
    }

    public Flux<ChatSession> getSessions() {
        return chatRepository.findAll();
    }

    public Mono<ChatSession> getSession(Mono<Long> id) {
        return chatRepository.findById(id);
    }

    public Mono<String> chat(ChatRequest request) {
        return getOrCreateConversation(Mono.justOrEmpty(request.conversationId()))
                .doOnNext(conversation -> saveMessage(conversation, request.prompt(), Author.USER))
                .flatMap(conversation -> {
                    List<Message> history = buildHistory(conversation);
                    return chatClient.prompt()
                            .messages(history)
                            .user(request.prompt())
                            .stream()
                            .content()
                            .collectList()
                            .map(strings -> String.join("", strings))
                            .doOnNext(content -> saveMessage(conversation, content, Author.LLM))
                            .flatMap(s -> chatRepository.save(conversation)
                                    .map(Iterable -> s));

                });
    }

    public Flux<String> streamChat(ChatRequest request) {
        return getOrCreateConversation(Mono.justOrEmpty(request.conversationId()))
            .doOnNext(conversation -> saveMessage(conversation, request.prompt(), Author.USER))
            .flatMapMany(conversation -> {
                        List<Message> history = buildHistory(conversation);
                        return chatClient.prompt()
                                .messages(history)
                                .user(request.prompt())
                                .stream()
                                .content();
                    });
//        StringBuilder aiContent = new StringBuilder();
//        return aiStream.doOnNext(aiContent::append)
//                .doOnComplete(() -> saveMessage(conv, aiContent.toString(), Owner.LLM))
//                .doOnError(e -> { /* Log error, rollback if needed */ });
    }

    private Mono<ChatSession> getOrCreateConversation(Mono<Long> id) {
        return chatRepository.findById(id)
                .switchIfEmpty(Mono.fromCallable(() -> new ChatSession("", new ArrayList<>())));
    }

    private void saveMessage(ChatSession conv, String content, Author role) {
        var message = new ChatMessage(content, role);
        conv.messages().add(message);
    }

    private List<Message> buildHistory(ChatSession conv) {
        return conv.messages()
                .stream()
                .map(msg -> msg.author() == Author.USER
                        ?  new UserMessage(msg.content())
                        : (Message) new AssistantMessage(msg.content()))
                .toList();
    }
}

