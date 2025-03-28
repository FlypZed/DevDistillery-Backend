package com.dev_disitllery.auth.controller;

import com.dev_disitllery.auth.model.CustomOAuth2User;
import com.dev_disitllery.auth.dto.ResponseDTO;
import com.dev_disitllery.auth.dto.UserDto;
import com.dev_disitllery.auth.model.User;
import com.dev_disitllery.auth.repository.UserRepository;
import com.dev_disitllery.auth.service.AuthService;
import com.dev_disitllery.auth.service.JwtService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.web.bind.annotation.*;
import jakarta.servlet.http.HttpServletResponse;

import java.io.IOException;
import java.util.Map;

@RestController
@RequestMapping("/api/auth")
@Tag(name = "Auth API", description = "API for authentication and user management")
public class AuthController {

    private final AuthService authService;
    private final JwtService jwtService;
    private final UserRepository userRepository;

    public AuthController(AuthService authService, JwtService jwtService, UserRepository userRepository) {
        this.authService = authService;
        this.jwtService = jwtService;
        this.userRepository = userRepository;
    }

    @GetMapping("/validate")
    @Operation(summary = "Validate JWT token", description = "Validate the JWT token provided in the Authorization header")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Token is valid"),
            @ApiResponse(responseCode = "401", description = "Token is invalid or expired")
    })
    public ResponseEntity<ResponseDTO<Boolean>> validateToken(@RequestHeader("Authorization") String token) {
        if (token != null && token.startsWith("Bearer ")) {
            String jwt = token.substring(7);
            if (jwtService.validateToken(jwt)) {
                return ResponseEntity.ok(ResponseDTO.ofSuccess(true));
            }
        }
        return ResponseEntity.status(401).body(ResponseDTO.ofError("Unauthorized", "Token is invalid or expired"));
    }

    @GetMapping("/user")
    @Operation(summary = "Get user info", description = "Get user information based on the JWT token provided in the Authorization header")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "User information retrieved successfully"),
            @ApiResponse(responseCode = "401", description = "Token is invalid or expired")
    })
    public ResponseEntity<?> getUserInfo(@RequestHeader("Authorization") String token) {
        if (token != null && token.startsWith("Bearer ")) {
            String jwt = token.substring(7);
            if (jwtService.validateToken(jwt)) {
                String email = jwtService.getEmailFromToken(jwt);
                User user = userRepository.findByEmail(email).orElse(null);
                if (user != null) {
                    UserDto userDto = new UserDto();
                    userDto.setEmail(user.getEmail());
                    userDto.setName(user.getName());
                    // Para futuras referencias aqui agregamos mas campos
                    return ResponseEntity.ok(ResponseDTO.ofSuccess(userDto));
                }
            }
        }
        return ResponseEntity.status(401).build();
    }

    @GetMapping("/oauth-user")
    @Operation(summary = "Get OAuth2 user info", description = "Get OAuth2 user information for the currently authenticated user")
    public ResponseEntity<?> getOAuthUserInfo(@AuthenticationPrincipal OAuth2User principal) {
        if (principal instanceof CustomOAuth2User) {
            CustomOAuth2User customUser = (CustomOAuth2User) principal;
            Map<String, Object> response = Map.of(
                    "email", customUser.getEmail(),
                    "name", customUser.getName(),
                    "picture", customUser.getPicture(),
                    "githubId", customUser.getGithubId(),
                    "githubLogin", customUser.getGithubLogin()
            );
            return ResponseEntity.ok(ResponseDTO.ofSuccess(response));
        }
        return ResponseEntity.status(401).build();
    }

    @PostMapping("/logout")
    @Operation(summary = "Logout", description = "Logout the user (token should be removed on the client side)")
    @ApiResponse(responseCode = "200", description = "Logout successful")
    public ResponseEntity<?> logout() {
        return ResponseEntity.ok().build();
    }

    @GetMapping("/token")
    @Operation(summary = "Get JWT token", description = "Get a JWT token for the authenticated OAuth2 user")
    @ApiResponses(value = {
            @ApiResponse(responseCode = "200", description = "Token generated successfully"),
            @ApiResponse(responseCode = "401", description = "User not authenticated")
    })
    public ResponseEntity<?> getToken(@AuthenticationPrincipal OAuth2User principal, HttpServletResponse response) throws IOException {
        if (principal != null) {
            String email = principal.getAttribute("email");
            String token = jwtService.generateToken(email);
            response.sendRedirect("http://localhost:5173/oauth-callback?token=" + token);
            return ResponseEntity.ok().build();
        }
        return ResponseEntity.status(401).build();
    }

    @GetMapping("/login-github")
    public ResponseEntity<?> loginWithGitHub(@RequestParam("code") String code) {
        authService.loginWithGitHub(code);
        return ResponseEntity.ok().build();
    }
}