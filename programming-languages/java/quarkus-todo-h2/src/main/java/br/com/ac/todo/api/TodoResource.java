package br.com.ac.todo.api;

import br.com.ac.todo.api.request.CreateTaskRequest;
import br.com.ac.todo.api.request.TaskRequest;
import br.com.ac.todo.api.response.TaskResponse;
import br.com.ac.todo.domain.Task;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.quarkus.hibernate.orm.panache.PanacheEntityBase;

import javax.transaction.Transactional;
import javax.validation.Valid;
import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.WebApplicationException;
import javax.ws.rs.core.Response;
import java.util.List;
import java.util.Optional;
import java.util.function.Function;

@Path("/tasks")
@Produces("application/json")
@Consumes("application/json")
public class TodoResource {

    private final ObjectMapper objectMapper = new ObjectMapper();

    public TodoResource() {
        objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
    }

    @GET
    public List<PanacheEntityBase> getAllTasks() {
        return Task.findAll().list();
    }

    @GET
    @Path("{id}")
    public TaskResponse getTaskById(@PathParam("id") Long id) {
        return mapToResponse()
                .apply(findById(id));
    }

    @POST
    @Transactional
    public TaskResponse create(@Valid  CreateTaskRequest taskRequest) {
        return mapToDto()
                .andThen(t -> {
                    t.persist();
                    return t;
                })
                .andThen(mapToResponse())
                .apply(taskRequest);
    }

    @DELETE
    @Path("{id}")
    @Transactional
    public void deleteById(@PathParam("id") Long id) {
        findById(id).delete();
    }

    @PUT
    @Path("{id}")
    @Transactional
    public TaskResponse updateById(@PathParam("id") Long id, TaskRequest taskRequest) {
        Task task = findById(id);
        task.setTitle(taskRequest.getTitle());
        task.setCompleted(taskRequest.isCompleted());
        task.persist();
        return mapToResponse()
                .apply(task);
    }

    private Task findById(Long id){
        return (Task) Optional.ofNullable(Task.findById(id))
                .orElseThrow(() -> new WebApplicationException("Not found!", Response.Status.NOT_FOUND));
    }

    private Function<Task, TaskResponse> mapToResponse() {
        return it -> objectMapper.convertValue(it, TaskResponse.class);
    }

    private Function<Object, Task> mapToDto() {
        return it -> objectMapper.convertValue(it, Task.class);
    }
}