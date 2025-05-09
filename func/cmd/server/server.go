package server

import (
	controllerBoard "func/internal/controller/board"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerUser "func/internal/controller/user"
	"log"

	repositoryBoard "func/internal/repository/board"
	repositoryProj "func/internal/repository/project"
	repositoryTask "func/internal/repository/task"
	repositoryUser "func/internal/repository/user"

	serviceBoard "func/internal/service/board"
	serviceProj "func/internal/service/project"
	serviceTask "func/internal/service/task"
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

	// Inicializa el hub de WebSocket
	wsHub := serviceWs.NewHub()
	go wsHub.Run()
	wsService := serviceWs.NewWebSocketService(wsHub)

	// Inicializa repositorios
	userRepo := repositoryUser.NewUserRepository(gormDB)
	projectRepo := repositoryProj.NewProjectRepository(sqlDB)
	taskRepo := repositoryTask.NewTaskRepository(gormDB)
	boardRepo := repositoryBoard.NewBoardRepository(gormDB)

	// Inicializa servicios
	userService := serviceUser.NewUserService(userRepo)
	projectService := serviceProj.NewProjectService(projectRepo)
	taskService := serviceTask.NewTaskService(taskRepo)
	boardService := serviceBoard.NewBoardService(boardRepo)

	// Inicializa controladores
	userController := controllerUser.NewUserController(userService)
	projectController := controllerProj.NewProjectController(projectService)
	taskController := controllerTask.NewTaskController(taskService, wsService)
	boardController := controllerBoard.NewBoardController(boardService, wsService)

	router := gin.Default()

	SetupRouter(
		router,
		userController,
		projectController,
		taskController,
		boardController,
	)

	return router
}
