package server

import (
	controllerBoard "func/internal/controller/board"
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"

	repositoryBoard "func/internal/repository/board"
	repositoryOrg "func/internal/repository/organization"
	repositoryProj "func/internal/repository/project"
	repositoryTask "func/internal/repository/task"
	repositoryTeam "func/internal/repository/team"
	repositoryUser "func/internal/repository/user"

	serviceBoard "func/internal/service/board"
	serviceOrg "func/internal/service/organization"
	serviceProj "func/internal/service/project"
	serviceTask "func/internal/service/task"
	serviceTeam "func/internal/service/team"
	serviceUser "func/internal/service/user"
	serviceWs "func/internal/service/websocket"

	"func/pkg/infrastructure"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	db := infrastructure.InitDB()

	// Inicializa el hub de WebSocket
	wsHub := serviceWs.NewHub()
	go wsHub.Run()
	wsService := serviceWs.NewWebSocketService(wsHub)

	// Inicializa repositorios
	userRepo := repositoryUser.NewUserRepository(db)
	teamRepo := repositoryTeam.NewTeamRepository(db)
	organizationRepo := repositoryOrg.NewOrganizationRepository(db)
	projectRepo := repositoryProj.NewProjectRepository(db)
	taskRepo := repositoryTask.NewTaskRepository(db)
	boardRepo := repositoryBoard.NewBoardRepository(db)

	// Inicializa servicios
	userService := serviceUser.NewUserService(userRepo)
	teamService := serviceTeam.NewTeamService(teamRepo)
	organizationService := serviceOrg.NewOrganizationService(organizationRepo)
	projectService := serviceProj.NewProjectService(projectRepo)
	taskService := serviceTask.NewTaskService(taskRepo)
	boardService := serviceBoard.NewBoardService(boardRepo)

	// Inicializa controladores pasando el wsService donde sea necesario
	userController := controllerUser.NewUserController(userService)
	teamController := controllerTeam.NewTeamController(teamService)
	organizationController := controllerOrg.NewOrganizationController(organizationService)
	projectController := controllerProj.NewProjectController(projectService)
	taskController := controllerTask.NewTaskController(taskService, wsService)
	boardController := controllerBoard.NewBoardController(boardService, wsService)

	router := gin.Default()

	SetupRouter(
		router,
		userController,
		teamController,
		organizationController,
		projectController,
		taskController,
		boardController,
	)

	return router
}
