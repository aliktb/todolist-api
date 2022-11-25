package com.example.todolistapi.controller;

import com.example.todolistapi.domain.User;
import com.example.todolistapi.identity.PlatformUser;
import com.example.todolistapi.service.UserService;
import liquibase.pro.packaged.G;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.web.bind.annotation.*;

@Slf4j
@CrossOrigin
@RestController
@RequestMapping("/users")
public class UserController {

    private UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @PostMapping("/add-user")
    public String addUser(@RequestBody User newUser) {

        userService.addNewUser(newUser);
        log.info("adding user");

        return "user added";
    }

    @GetMapping("/get-user")
    public ResponseEntity<User> fetchUser(@RequestParam String userId) {

        User result = userService.getUser(userId);

        return ResponseEntity.ok(result);
    }

    @GetMapping("/sample")
    public String sample() {

        JwtAuthenticationToken auth = (JwtAuthenticationToken) SecurityContextHolder.getContext().getAuthentication();

        PlatformUser currentUser = new PlatformUser(auth);

        return currentUser.getSubject();
    }

    // checks if users exists in DB
    @GetMapping("/exist")
    public boolean doesUserExist() {

        JwtAuthenticationToken auth = (JwtAuthenticationToken) SecurityContextHolder.getContext().getAuthentication();

        PlatformUser currentUser = new PlatformUser(auth);

//        User newUser = userService.getUser(currentUser.getSubject());
//
        return userService.doesUserExistinDB(currentUser);


    }
}