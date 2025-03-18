package com.dev_disitllery.auth.exception;

public class OAuth2AuthenticationProcessingException extends RuntimeException {
    public OAuth2AuthenticationProcessingException(String message, Throwable cause) {
        super(message, cause);
    }
}