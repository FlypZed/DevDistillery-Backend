package server

import (
	controllerBoard "func/internal/controller/board"
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

		teamRoutes := api.Group("/teams")
		{
			teamRoutes.POST("/", teamController.Create)
			teamRoutes.GET("/:id", teamController.Get)
			teamRoutes.GET("/organization/:orgId", teamController.GetByOrganization)
			teamRoutes.PUT("/:id", teamController.Update)
			teamRoutes.DELETE("/:id", teamController.Delete)
		}

		organizationRoutes := api.Group("/organizations")
		{
			organizationRoutes.POST("/", organizationController.Create)
			organizationRoutes.GET("/", organizationController.GetAll)
			organizationRoutes.GET("/:id", organizationController.Get)
			organizationRoutes.PUT("/:id", organizationController.Update)
			organizationRoutes.DELETE("/:id", organizationController.Delete)
		}

		projectRoutes := api.Group("/projects")
		{
			projectRoutes.POST("/", projectController.Create)
			projectRoutes.GET("/:id", projectController.Get)
			projectRoutes.GET("/user/:userId", projectController.GetByUser)
			projectRoutes.GET("/organization/:orgId", projectController.GetByOrganization)
			projectRoutes.PUT("/:id", projectController.Update)
			projectRoutes.PUT("/:id/team", projectController.AssignTeam)
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
