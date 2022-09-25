package com.example.todolistapi.controller;

import com.example.todolistapi.domain.Task;
import com.example.todolistapi.service.TaskService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@Slf4j
@CrossOrigin
@RestController
@RequestMapping("/tasks")
public class TaskController {

    private TaskService service;

    public TaskController(TaskService service) {
        this.service = service;
    }

    @GetMapping("/getAllTasks")
    public ResponseEntity<List<Task>> getAllTasks(){

        return new ResponseEntity<>(service.getAllTasks(), HttpStatus.OK);
    }

    @PostMapping("/addNewTask")
    public ResponseEntity<String> addNewTask(@RequestBody Task newTask){

        return new ResponseEntity<>(service.addNewTask(newTask), HttpStatus.OK);
    }
}
