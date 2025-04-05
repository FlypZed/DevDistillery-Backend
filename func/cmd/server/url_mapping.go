package server

import (
	controllerBoard "func/internal/controller/board"
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	router *gin.Engine,
	userController *controllerUser.UserController,
	teamController *controllerTeam.TeamController,
	organizationController *controllerOrg.OrganizationController,
	projectController *controllerProj.ProjectController,
	taskController *controllerTask.TaskController,
	boardController *controllerBoard.BoardController,
) {
	api := router.Group("/api")
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("/", userController.CreateUser)
			userRoutes.GET("/:id", userController.GetUser)
			userRoutes.PUT("/:id", userController.UpdateUser)
			userRoutes.DELETE("/:id", userController.DeleteUser)
		}

		teamRoutes := api.Group("/teams")
		{
			teamRoutes.POST("/", teamController.CreateTeam)
			teamRoutes.GET("/:id", teamController.GetTeam)
			teamRoutes.PUT("/:id", teamController.UpdateTeam)
			teamRoutes.DELETE("/:id", teamController.DeleteTeam)
		}

		organizationRoutes := api.Group("/organizations")
		{
			organizationRoutes.POST("/", organizationController.CreateOrganization)
			organizationRoutes.GET("/:id", organizationController.GetOrganization)
			organizationRoutes.PUT("/:id", organizationController.UpdateOrganization)
			organizationRoutes.DELETE("/:id", organizationController.DeleteOrganization)
		}

		projectRoutes := api.Group("/projects")
		{
			projectRoutes.POST("/", projectController.CreateProject)
			projectRoutes.GET("/:id", projectController.GetProject)
			projectRoutes.PUT("/:id", projectController.UpdateProject)
			projectRoutes.DELETE("/:id", projectController.DeleteProject)
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
