package com.ac.chat.service;

import com.ac.chat.domain.Author;
import com.ac.chat.domain.ChatMessage;
import com.ac.chat.domain.ChatRequest;
import com.ac.chat.domain.ChatSession;
import com.ac.chat.repository.ChatRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
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
    private static final Logger log = LoggerFactory.getLogger(ChatService.class);

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
                .flatMap(conversation -> {
                    List<Message> history = buildHistory(conversation);
                    saveMessage(conversation, request.prompt(), Author.USER);
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
            .flatMapMany(conversation -> {
                List<Message> history = buildHistory(conversation);
                StringBuilder aiContent = new StringBuilder();

                saveMessage(conversation, request.prompt(), Author.USER);

                return chatClient.prompt()
                        .messages(history)
                        .user(request.prompt())
                        .stream()
                        .content()
                        .doOnNext(chunk -> {
                            aiContent.append(chunk);
                            log.info("LLM stream chunk conversationId={} chunk=\"{}\"", conversation.id(), printable(chunk));
                        })
                        .doFinally(signalType -> {
                            if (!aiContent.isEmpty()) {
                                saveMessage(conversation, aiContent.toString(), Author.LLM);
                            }
                            chatRepository.save(conversation);
                            log.info("Finished LLM stream conversationId={} signal={} responseLength={}",
                                    conversation.id(), signalType, aiContent.length());
                        });
            });
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

    private String printable(String value) {
        if (value == null) {
            return "";
        }

        return value
                .replace("\r", "\\r")
                .replace("\n", "\\n");
    }
}
