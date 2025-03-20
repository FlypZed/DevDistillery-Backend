package com.dev_disitllery.auth.dto;

public class ResponseDTO<T> {

    private static final String OPERATION_SUCCESSFUL = "Operation successful";
    private String message;
    private String errorCode;
    private boolean success;
    private T content;

    public ResponseDTO(String message, String errorCode, boolean success, T content) {
        this.message = message;
        this.errorCode = errorCode;
        this.success = success;
        this.content = content;
    }

    private ResponseDTO(String message, T content) {
        this(message, null, true, content);
    }

    private ResponseDTO(String message, String errorCode) {
        this(message, errorCode, false, null);
    }

    public static <T> ResponseDTO<T> ofError(String reason, String message) {
        return new ResponseDTO<>(message, reason);
    }

    public static <T> ResponseDTO<T> ofSuccess(T content) {
        return new ResponseDTO<>(OPERATION_SUCCESSFUL, content);
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public void setErrorCode(String errorCode) {
        this.errorCode = errorCode;
    }

    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }

    public T getContent() {
        return content;
    }

    public void setContent(T content) {
        this.content = content;
    }
}
