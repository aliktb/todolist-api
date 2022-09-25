package com.example.todolistapi.service;

import com.example.todolistapi.domain.Task;
import com.example.todolistapi.repository.TaskRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class TaskService {

    @Autowired
    private final TaskRepository taskRepository;

    public TaskService(TaskRepository taskRepository) {
        this.taskRepository = taskRepository;
    }

    public List<Task> getAllTasks(){
        return taskRepository.findAll();
    }

    public String addNewTask(Task newTask){

        taskRepository.saveAndFlush(newTask);
        return "added";
    }

}
