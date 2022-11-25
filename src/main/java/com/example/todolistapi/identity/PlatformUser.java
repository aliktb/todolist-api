package com.example.todolistapi.identity;

import lombok.Getter;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;

@Getter
public class PlatformUser {

    private final String subject;

    public PlatformUser(JwtAuthenticationToken jwtAuthToken) {
        this.subject = jwtAuthToken.getName();
    }
}
