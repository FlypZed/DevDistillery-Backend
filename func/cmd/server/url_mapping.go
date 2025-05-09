package server

import (
	controllerBoard "func/internal/controller/board"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerUser "func/internal/controller/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	router *gin.Engine,
	userController *controllerUser.UserController,
	projectController *controllerProj.ProjectController,
	taskController *controllerTask.TaskController,
	boardController *controllerBoard.BoardController,
) {
	router.Use(cors.Default())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("/", userController.CreateUser)
			userRoutes.GET("/:id", userController.GetUser)
			userRoutes.PUT("/:id", userController.UpdateUser)
			userRoutes.DELETE("/:id", userController.DeleteUser)
		}

		projectRoutes := api.Group("/projects")
		{
			projectRoutes.POST("/", projectController.Create)
			projectRoutes.GET("/:id", projectController.Get)
			projectRoutes.GET("/user/:userId", projectController.GetByUser)
			projectRoutes.PUT("/:id", projectController.Update)
			projectRoutes.DELETE("/:id", projectController.Delete)
		}

		taskRoutes := api.Group("/tasks")
		{
			taskRoutes.POST("/", taskController.CreateTask)
			taskRoutes.GET("/:id", taskController.GetTask)
			taskRoutes.PUT("/:id", taskController.UpdateTask)
			taskRoutes.DELETE("/:id", taskController.DeleteTask)
		}

		boardRoutes := api.Group("/boards")
		{
			boardRoutes.GET("/:id", boardController.GetBoard)
			boardRoutes.PUT("/:id", boardController.SaveBoard)
			boardRoutes.GET("/ws/:boardId", boardController.HandleWebSocket)
		}
	}
}
