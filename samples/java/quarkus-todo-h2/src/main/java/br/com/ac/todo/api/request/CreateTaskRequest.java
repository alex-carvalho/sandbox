package br.com.ac.todo.api.request;

import javax.validation.constraints.NotBlank;

public class CreateTaskRequest {

    @NotBlank(message="Title not be blank.")
    private String title;

    public CreateTaskRequest() {
    }

    public CreateTaskRequest(@NotBlank(message = "Title not be blank.") String title) {
        this.title = title;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }
}
