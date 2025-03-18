package com.dev_disitllery.auth.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@Controller
public class AuthController {

    @GetMapping("/login")
    public String login(@RequestParam(value = "error", required = false) String error) {
        if (error != null) {
            return "login-error";
        }
        return "login";
    }

    @GetMapping("/home")
    public String home() {
        return "home";
    }
}