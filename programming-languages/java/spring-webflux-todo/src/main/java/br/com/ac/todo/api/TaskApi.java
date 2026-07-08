package br.com.ac.todo.api;


import br.com.ac.todo.api.request.CreateTaskRequest;
import br.com.ac.todo.api.request.TaskRequest;
import br.com.ac.todo.api.response.TaskResponse;
import br.com.ac.todo.dto.Task;
import br.com.ac.todo.repository.TaskRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.function.Function;

@RestController
@RequestMapping("/tasks")
public class TaskApi {

    private final TaskRepository taskRepository;
    private final ObjectMapper objectMapper;

    public TaskApi(TaskRepository taskRepository, ObjectMapper objectMapper) {
        this.taskRepository = taskRepository;
        this.objectMapper = objectMapper;
    }

    @GetMapping
    public Flux<TaskResponse> getAllTasks() {
        return taskRepository.findAll()
                .map(mapToResponse());
    }

    @GetMapping("/{id}")
    public Mono<TaskResponse> getTaskById(@PathVariable("id") String id) {
        return taskRepository.findTaskById(id)
                .map(mapToResponse());
    }

    @PostMapping
    public Mono<TaskResponse> create(@RequestBody CreateTaskRequest taskRequest) {
        return Mono.just(taskRequest)
                .map(it -> objectMapper.convertValue(it, Task.class))
                .flatMap(taskRepository::create)
                .map(mapToResponse());
    }

    @DeleteMapping("/{id}")
    public Mono<ResponseEntity> deleteById(@PathVariable String id) {
        return taskRepository.deleteById(id)
                .flatMap(aVoid -> Mono.just(new ResponseEntity(HttpStatus.OK)));
    }

    @PutMapping("/{id}")
    public Mono<TaskResponse> updateById(@PathVariable String id, @RequestBody TaskRequest taskRequest) {
        return Mono.just(taskRequest)
                .map(mapToDto())
                .flatMap(it -> taskRepository.updateById(id, it))
                .map(mapToResponse());
    }

    private Function<TaskRequest, Task> mapToDto() {
        return it -> objectMapper.convertValue(it, Task.class);
    }


    private Function<Task, TaskResponse> mapToResponse() {
        return it -> objectMapper.convertValue(it, TaskResponse.class);
    }
}