package response

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Content interface{} `json:"content,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Success(c *gin.Context, status int, content interface{}, message string) {
	c.JSON(status, Response{
		Success: true,
		Content: content,
		Message: message,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}
