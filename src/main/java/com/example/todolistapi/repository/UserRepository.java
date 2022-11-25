package com.example.todolistapi.repository;

import com.example.todolistapi.domain.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface UserRepository  extends JpaRepository<User, Integer> {

    // SQL Query
    @Query(value = "SELECT * FROM aquarius.user WHERE user_id = ?1", nativeQuery = true)
    List<User> findUserWithId(String userId);
}
