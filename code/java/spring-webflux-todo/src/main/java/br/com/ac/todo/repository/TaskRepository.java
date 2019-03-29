package br.com.ac.todo.repository;


import br.com.ac.todo.dto.Task;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Component
public class TaskRepository {

    private final Map<String, Task> tasks = new HashMap<>();


    public Flux<Task> findAll() {
        return Flux.fromIterable(tasks.values());
    }


    public Mono<Task> findTaskById(String id) {
        return Mono.justOrEmpty(tasks.get(id));
    }

    public Mono<Task> create(Task task) {
        task.setId(UUID.randomUUID().toString());
        tasks.put(task.getId(), task);
        return Mono.just(task);
    }

    public Mono<Void> deleteById(String id) {
        tasks.remove(id);
        return Mono.empty();
    }

    public Mono<Task> updateById(String id, Task task) {
        task.setId(id);
        tasks.put(id, task);
        return Mono.just(task);
    }
}
