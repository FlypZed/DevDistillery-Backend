package controller

import (
	"func/internal/domain"
	service "func/internal/service/user"
	"func/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user data: "+err.Error())
		return
	}

	if err := uc.userService.CreateUser(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, user, "User created successfully")
}

func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.userService.GetUser(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, user, "User retrieved successfully")
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid user data: "+err.Error())
		return
	}

	user.ID = id
	if err := uc.userService.UpdateUser(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, user, "User updated successfully")
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.userService.DeleteUser(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil, "User deleted successfully")
}
