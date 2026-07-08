package br.com.ac.todo;

import br.com.ac.todo.api.request.CreateTaskRequest;
import br.com.ac.todo.api.request.TaskRequest;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.quarkus.test.junit.QuarkusTest;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.Test;

import javax.ws.rs.core.Response;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.notNullValue;

@QuarkusTest
class TodoResourceTest {

    @Test
    void testGetTasks() {
        given()
                .when().get("/tasks")
                .then()
                .statusCode(200)
                .body(is("[]"));
    }

    @Test
    void testCreateTask() throws JsonProcessingException {
        String payload = new ObjectMapper().writeValueAsString(new CreateTaskRequest("teste"));
        given()
                .contentType(ContentType.JSON)
                .body(payload)
                .when().post("/tasks")
                .then()
                .statusCode(Response.Status.OK.getStatusCode())
                .body("id", notNullValue())
                .body("completed", equalTo(false));
    }

    @Test
    void testUpdateTask() throws JsonProcessingException {
        String payloadCreate = new ObjectMapper().writeValueAsString(new CreateTaskRequest("teste"));
        given()
                .contentType(ContentType.JSON)
                .body(payloadCreate)
                .when().post("/tasks");

        TaskRequest taskRequest = new TaskRequest("teste2", true);
        String payloadUpdate = new ObjectMapper().writeValueAsString(taskRequest);

        given()
                .contentType(ContentType.JSON)
                .body(payloadUpdate)
                .when().put("/tasks/{id}", 1)
                .then()
                .statusCode(Response.Status.OK.getStatusCode())
                .body("id", notNullValue())
                .body("title", equalTo(taskRequest.getTitle()))
                .body("completed", equalTo(taskRequest.isCompleted()));
    }

}
