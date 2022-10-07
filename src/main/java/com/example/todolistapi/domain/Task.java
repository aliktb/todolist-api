package com.example.todolistapi.domain;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.RequiredArgsConstructor;

import javax.persistence.Entity;
import javax.persistence.Id;
import java.util.Date;

@Entity
@Data
@AllArgsConstructor
@RequiredArgsConstructor
public class Task {

    @Id
    private int taskId;

    private String taskName;
    private String taskDescription;
    private String owner;
    private Date dueDate;
    private Date creationDate;

}
