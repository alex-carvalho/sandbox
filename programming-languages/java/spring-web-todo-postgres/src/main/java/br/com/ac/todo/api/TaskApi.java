package br.com.ac.todo.api;


import br.com.ac.todo.api.request.CreateTaskRequest;
import br.com.ac.todo.api.request.TaskRequest;
import br.com.ac.todo.api.response.TaskResponse;
import br.com.ac.todo.domain.Task;
import br.com.ac.todo.repository.TaskRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.function.Function;
import java.util.stream.Collectors;

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
    public List<TaskResponse> getAllTasks() {
        return taskRepository.findAll()
                .stream()
                .map(task -> objectMapper.convertValue(task, TaskResponse.class))
                .collect(Collectors.toList());
    }

    @GetMapping("/{id}")
    public TaskResponse getTaskById(@PathVariable("id") Long id) {
        return mapToResponse()
                .apply(taskRepository.findById(id).orElseThrow(() -> new RuntimeException("Not found!")));
    }

    @PostMapping
    public TaskResponse create(@RequestBody CreateTaskRequest taskRequest) {
        return mapToDto()
                .andThen(taskRepository::save)
                .andThen(mapToResponse())
                .apply(taskRequest);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity deleteById(@PathVariable Long id) {
        taskRepository.deleteById(id);
        return ResponseEntity.ok().build();
    }

    @PutMapping("/{id}")
    public TaskResponse updateById(@PathVariable Long id, @RequestBody TaskRequest taskRequest) {
        taskRequest.setId(id);
        return mapToDto()
                .andThen(taskRepository::save)
                .andThen(mapToResponse())
                .apply(taskRequest);
    }

    private Function<Object, Task> mapToDto() {
        return it -> objectMapper.convertValue(it, Task.class);
    }


    private Function<Task, TaskResponse> mapToResponse() {
        return it -> objectMapper.convertValue(it, TaskResponse.class);
    }
}