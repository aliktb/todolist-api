package com.example.todolistapi.Domain;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.RequiredArgsConstructor;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Table;
import java.util.Date;

@Entity
@Data
@AllArgsConstructor
@RequiredArgsConstructor
@Table(name = "tasks")
public class Task {

    @Id
    private int taskId;

    private String taskName;
    private String taskDescription;
    private String owner;
    private Date dueDate;
    private Date creationDate;

}
