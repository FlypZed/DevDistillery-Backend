package server

import (
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"

	repositoryOrg "func/internal/repository/organization"
	repositoryProj "func/internal/repository/project"
	repositoryTask "func/internal/repository/task"
	repositoryTeam "func/internal/repository/team"
	repositoryUser "func/internal/repository/user"

	serviceOrg "func/internal/service/organization"
	serviceProj "func/internal/service/project"
	serviceTask "func/internal/service/task"
	serviceTeam "func/internal/service/team"
	serviceUser "func/internal/service/user"

	"func/pkg/infrastructure"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	db := infrastructure.InitDB()

	userRepo := repositoryUser.NewUserRepository(db)
	teamRepo := repositoryTeam.NewTeamRepository(db)
	organizationRepo := repositoryOrg.NewOrganizationRepository(db)
	projectRepo := repositoryProj.NewProjectRepository(db)
	taskRepo := repositoryTask.NewTaskRepository(db)

	userService := serviceUser.NewUserService(userRepo)
	teamService := serviceTeam.NewTeamService(teamRepo)
	organizationService := serviceOrg.NewOrganizationService(organizationRepo)
	projectService := serviceProj.NewProjectService(projectRepo)
	taskService := serviceTask.NewTaskService(taskRepo)

	userController := controllerUser.NewUserController(userService)
	teamController := controllerTeam.NewTeamController(teamService)
	organizationController := controllerOrg.NewOrganizationController(organizationService)
	projectController := controllerProj.NewProjectController(projectService)
	taskController := controllerTask.NewTaskController(taskService)

	router := gin.Default()

	SetupRouter(
		router,
		userController,
		teamController,
		organizationController,
		projectController,
		taskController,
	)

	router.GET("/ws", func(c *gin.Context) {
		infrastructure.WebSocketHandler(c.Writer, c.Request)
	})

	return router
}
