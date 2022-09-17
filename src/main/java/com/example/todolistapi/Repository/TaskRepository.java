package com.example.todolistapi.Repository;

import com.example.todolistapi.Domain.Task;
import org.springframework.data.jpa.repository.JpaRepository;

public interface TaskRepository extends JpaRepository<Task, Integer> {

}
