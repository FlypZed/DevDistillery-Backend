package server

import (
	controllerBoard "func/internal/controller/board"
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title DevDistillery API
// @version 1.0
// @description This is the API documentation for DevDistillery.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

func SetupRouter(
	router *gin.Engine,
	userController *controllerUser.UserController,
	teamController *controllerTeam.TeamController,
	organizationController *controllerOrg.OrganizationController,
	projectController *controllerProj.ProjectController,
	taskController *controllerTask.TaskController,
	boardController *controllerBoard.BoardController,
) {
	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		userRoutes := api.Group("/users")
		{
			// @Summary Create a new user
			// @Description Create a new user with the input payload
			// @Tags users
			// @Accept  json
			// @Produce  json
			// @Param   user      body     controllerUser.CreateUserRequest     true  "User payload"
			// @Success 200 {object} controllerUser.UserResponse
			// @Router /users [post]
			userRoutes.POST("/", userController.CreateUser)
			// @Summary Get a user by ID
			// @Description Get details of a user by ID
			// @Tags users
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "User ID"
			// @Success 200 {object} controllerUser.UserResponse
			// @Router /users/{id} [get]
			userRoutes.GET("/:id", userController.GetUser)
			// @Summary Update a user by ID
			// @Description Update details of a user by ID
			// @Tags users
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "User ID"
			// @Param   user      body     controllerUser.UpdateUserRequest     true  "User payload"
			// @Success 200 {object} controllerUser.UserResponse
			// @Router /users/{id} [put]
			userRoutes.PUT("/:id", userController.UpdateUser)
			// @Summary Delete a user by ID
			// @Description Delete a user by ID
			// @Tags users
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "User ID"
			// @Success 200 {object} controllerUser.UserResponse
			// @Router /users/{id} [delete]
			userRoutes.DELETE("/:id", userController.DeleteUser)
		}

		teamRoutes := api.Group("/teams")
		{
			// @Summary Create a new team
			// @Description Create a new team with the input payload
			// @Tags teams
			// @Accept  json
			// @Produce  json
			// @Param   team      body     controllerTeam.CreateTeamRequest     true  "Team payload"
			// @Success 200 {object} controllerTeam.TeamResponse
			// @Router /teams [post]
			teamRoutes.POST("/", teamController.CreateTeam)
			// @Summary Get a team by ID
			// @Description Get details of a team by ID
			// @Tags teams
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Team ID"
			// @Success 200 {object} controllerTeam.TeamResponse
			// @Router /teams/{id} [get]
			teamRoutes.GET("/:id", teamController.GetTeam)
			// @Summary Update a team by ID
			// @Description Update details of a team by ID
			// @Tags teams
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Team ID"
			// @Param   team      body     controllerTeam.UpdateTeamRequest     true  "Team payload"
			// @Success 200 {object} controllerTeam.TeamResponse
			// @Router /teams/{id} [put]
			teamRoutes.PUT("/:id", teamController.UpdateTeam)
			// @Summary Delete a team by ID
			// @Description Delete a team by ID
			// @Tags teams
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Team ID"
			// @Success 200 {object} controllerTeam.TeamResponse
			// @Router /teams/{id} [delete]
			teamRoutes.DELETE("/:id", teamController.DeleteTeam)
		}

		organizationRoutes := api.Group("/organizations")
		{
			// @Summary Create a new organization
			// @Description Create a new organization with the input payload
			// @Tags organizations
			// @Accept  json
			// @Produce  json
			// @Param   organization      body     controllerOrg.CreateOrganizationRequest     true  "Organization payload"
			// @Success 200 {object} controllerOrg.OrganizationResponse
			// @Router /organizations [post]
			organizationRoutes.POST("/", organizationController.CreateOrganization)
			// @Summary Get an organization by ID
			// @Description Get details of an organization by ID
			// @Tags organizations
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Organization ID"
			// @Success 200 {object} controllerOrg.OrganizationResponse
			// @Router /organizations/{id} [get]
			organizationRoutes.GET("/:id", organizationController.GetOrganization)
			// @Summary Update an organization by ID
			// @Description Update details of an organization by ID
			// @Tags organizations
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Organization ID"
			// @Param   organization      body     controllerOrg.UpdateOrganizationRequest     true  "Organization payload"
			// @Success 200 {object} controllerOrg.OrganizationResponse
			// @Router /organizations/{id} [put]
			organizationRoutes.PUT("/:id", organizationController.UpdateOrganization)
			// @Summary Delete an organization by ID
			// @Description Delete an organization by ID
			// @Tags organizations
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Organization ID"
			// @Success 200 {object} controllerOrg.OrganizationResponse
			// @Router /organizations/{id} [delete]
			organizationRoutes.DELETE("/:id", organizationController.DeleteOrganization)
		}

		projectRoutes := api.Group("/projects")
		{
			// @Summary Create a new project
			// @Description Create a new project with the input payload
			// @Tags projects
			// @Accept  json
			// @Produce  json
			// @Param   project      body     controllerProj.CreateProjectRequest     true  "Project payload"
			// @Success 200 {object} controllerProj.ProjectResponse
			// @Router /projects [post]
			projectRoutes.POST("/", projectController.CreateProject)
			// @Summary Get a project by ID
			// @Description Get details of a project by ID
			// @Tags projects
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Project ID"
			// @Success 200 {object} controllerProj.ProjectResponse
			// @Router /projects/{id} [get]
			projectRoutes.GET("/:id", projectController.GetProject)
			// @Summary Update a project by ID
			// @Description Update details of a project by ID
			// @Tags projects
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Project ID"
			// @Param   project      body     controllerProj.UpdateProjectRequest     true  "Project payload"
			// @Success 200 {object} controllerProj.ProjectResponse
			// @Router /projects/{id} [put]
			projectRoutes.PUT("/:id", projectController.UpdateProject)
			// @Summary Delete a project by ID
			// @Description Delete a project by ID
			// @Tags projects
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Project ID"
			// @Success 200 {object} controllerProj.ProjectResponse
			// @Router /projects/{id} [delete]
			projectRoutes.DELETE("/:id", projectController.DeleteProject)
		}

		taskRoutes := api.Group("/tasks")
		{
			// @Summary Create a new task
			// @Description Create a new task with the input payload
			// @Tags tasks
			// @Accept  json
			// @Produce  json
			// @Param   task      body     controllerTask.CreateTaskRequest     true  "Task payload"
			// @Success 200 {object} controllerTask.TaskResponse
			// @Router /tasks [post]
			taskRoutes.POST("/", taskController.CreateTask)
			// @Summary Get a task by ID
			// @Description Get details of a task by ID
			// @Tags tasks
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Task ID"
			// @Success 200 {object} controllerTask.TaskResponse
			// @Router /tasks/{id} [get]
			taskRoutes.GET("/:id", taskController.GetTask)
			// @Summary Update a task by ID
			// @Description Update details of a task by ID
			// @Tags tasks
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Task ID"
			// @Param   task      body     controllerTask.UpdateTaskRequest     true  "Task payload"
			// @Success 200 {object} controllerTask.TaskResponse
			// @Router /tasks/{id} [put]
			taskRoutes.PUT("/:id", taskController.UpdateTask)
			// @Summary Delete a task by ID
			// @Description Delete a task by ID
			// @Tags tasks
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Task ID"
			// @Success 200 {object} controllerTask.TaskResponse
			// @Router /tasks/{id} [delete]
			taskRoutes.DELETE("/:id", taskController.DeleteTask)
		}

		boardRoutes := api.Group("/boards")
		{
			// @Summary Get a board by ID
			// @Description Get details of a board by ID
			// @Tags boards
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Board ID"
			// @Success 200 {object} controllerBoard.BoardResponse
			// @Router /boards/{id} [get]
			boardRoutes.GET("/:id", boardController.GetBoard)
			// @Summary Save a board by ID
			// @Description Save details of a board by ID
			// @Tags boards
			// @Accept  json
			// @Produce  json
			// @Param   id      path     string     true  "Board ID"
			// @Param   board      body     controllerBoard.SaveBoardRequest     true  "Board payload"
			// @Success 200 {object} controllerBoard.BoardResponse
			// @Router /boards/{id} [put]
			boardRoutes.PUT("/:id", boardController.SaveBoard)
			// @Summary Handle WebSocket connection for a board
			// @Description Handle WebSocket connection for a board by ID
			// @Tags boards
			// @Accept  json
			// @Produce  json
			// @Param   boardId      path     string     true  "Board ID"
			// @Success 200 {object} controllerBoard.BoardResponse
			// @Router /boards/ws/{boardId} [get]
			boardRoutes.GET("/ws/:boardId", boardController.HandleWebSocket)
		}
	}
}
