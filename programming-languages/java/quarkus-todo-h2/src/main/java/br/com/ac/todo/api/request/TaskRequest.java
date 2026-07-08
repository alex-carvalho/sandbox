package br.com.ac.todo.api.request;

public class TaskRequest {

    private String title;
    private boolean completed;


    public TaskRequest() {
    }

    public TaskRequest(String title, boolean completed) {
        this.title = title;
        this.completed = completed;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public boolean isCompleted() {
        return completed;
    }

    public void setCompleted(boolean completed) {
        this.completed = completed;
    }
}
