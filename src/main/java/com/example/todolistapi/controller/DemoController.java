package com.example.todolistapi.controller;

import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@Slf4j
@RestController
@RequestMapping("/demo")
public class DemoController {

    @GetMapping("/test")
    public ResponseEntity<String> testEndpoint(){

        log.info("Fetching demo/test endpoint");
        return new ResponseEntity<String>("Running!", HttpStatus.OK);
    }

}
