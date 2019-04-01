package br.com.ac.todo.repository;


import br.com.ac.todo.domain.Task;
import org.springframework.data.jpa.repository.JpaRepository;

public interface TaskRepository extends JpaRepository<Task, Long> {
}
