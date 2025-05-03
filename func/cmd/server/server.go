package server

import (
	controllerBoard "func/internal/controller/board"
	controllerOrg "func/internal/controller/organization"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerTeam "func/internal/controller/team"
	controllerUser "func/internal/controller/user"
	"log"

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
	gormDB := infrastructure.InitDB()

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB", err)
	}
	teamRepo := repositoryTeam.NewTeamRepository(sqlDB)

	// Inicializa el hub de WebSocket
	wsHub := serviceWs.NewHub()
	go wsHub.Run()
	wsService := serviceWs.NewWebSocketService(wsHub)

	// Inicializa repositorios
	userRepo := repositoryUser.NewUserRepository(gormDB)
	organizationRepo := repositoryOrg.NewOrganizationRepository(sqlDB)
	projectRepo := repositoryProj.NewProjectRepository(sqlDB)
	taskRepo := repositoryTask.NewTaskRepository(gormDB)
	boardRepo := repositoryBoard.NewBoardRepository(gormDB)

	// Inicializa servicios
	userService := serviceUser.NewUserService(userRepo)
	teamService := serviceTeam.NewTeamService(teamRepo)
	organizationService := serviceOrg.NewOrganizationService(organizationRepo)
	projectService := serviceProj.NewProjectService(projectRepo)
	taskService := serviceTask.NewTaskService(taskRepo)
	boardService := serviceBoard.NewBoardService(boardRepo)

	// Inicializa controladores
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
