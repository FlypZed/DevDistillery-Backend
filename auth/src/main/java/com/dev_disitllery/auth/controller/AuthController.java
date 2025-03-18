package com.dev_disitllery.auth.controller;

import com.dev_disitllery.auth.model.User;
import com.dev_disitllery.auth.repository.UserRepository;
import com.dev_disitllery.auth.service.JwtService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/auth")
public class AuthController {

    private final JwtService jwtService;
    private final UserRepository userRepository;

    public AuthController(JwtService jwtService, UserRepository userRepository) {
        this.jwtService = jwtService;
        this.userRepository = userRepository;
    }

    @GetMapping("/validate")
    public ResponseEntity<?> validateToken(@RequestHeader("Authorization") String token) {
        if (token != null && token.startsWith("Bearer ")) {
            String jwt = token.substring(7);
            if (jwtService.validateToken(jwt)) {
                return ResponseEntity.ok().build();
            }
        }
        return ResponseEntity.status(401).build();
    }

    @GetMapping("/user")
    public ResponseEntity<?> getUserInfo(@RequestHeader("Authorization") String token) {
        if (token != null && token.startsWith("Bearer ")) {
            String jwt = token.substring(7);
            if (jwtService.validateToken(jwt)) {
                String email = jwtService.getEmailFromToken(jwt);
                User user = userRepository.findByEmail(email).orElse(null);
                if (user != null) {
                    return ResponseEntity.ok(user);
                }
            }
        }
        return ResponseEntity.status(401).build();
    }

    @PostMapping("/logout")
    public ResponseEntity<?> logout() {
        // El logout debe ser manejado en el frontend eliminando el token.
        return ResponseEntity.ok().build();
    }

    @GetMapping("/token")
    public ResponseEntity<?> getToken(@AuthenticationPrincipal OAuth2User principal) {
        if (principal != null) {
            String email = principal.getAttribute("email");
            String token = jwtService.generateToken(email);
            return ResponseEntity.ok(token);
        }
        return ResponseEntity.status(401).build();
    }
}