package server

import (
	controllerBoard "func/internal/controller/board"
	controllerProj "func/internal/controller/project"
	controllerTask "func/internal/controller/task"
	controllerUser "func/internal/controller/user"
	"log"
	"net/http"
	"time"

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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gormDB := infrastructure.InitDB()

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("failed to get underlying sql.DB", err)
	}

	// Configuración inicial del router
	router := gin.New()

	// Middleware de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.Use(gin.Recovery())

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

	SetupRouter(
		router,
		userController,
		projectController,
		taskController,
		boardController,
	)

	return router
}
