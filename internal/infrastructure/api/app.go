package api

import (
	"context"
	"database/sql"
	"fmt"
	"lab1_isbn/internal/core/ports/repositories"
	"lab1_isbn/internal/core/record"
	"lab1_isbn/internal/infrastructure/api/handler"
	"lab1_isbn/internal/infrastructure/observability"
	"lab1_isbn/internal/infrastructure/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type dependencies struct {
	RecordHandler *handler.RecordHandler
	OtelShutdown  func(context.Context) error
}

type infrastructures struct {
	Storage repositories.TaskRepository
}

type App struct {
	Server *gin.Engine
	deps   *dependencies
	infra  *infrastructures
}

func NewApp() *App {
	a := &App{}
	a.setupInfrastructure()
	a.setupDependencies()
	a.setupServer()

	return a
}

func (a *App) setupInfrastructure() {
	dsn := "root:123@tcp(localhost:3306)/isbn"
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db := bun.NewDB(sqlDB, mysqldialect.New())
	storage := storage.NewStorage(db)

	a.infra = &infrastructures{
		Storage: storage,
	}
}

func (a *App) setupDependencies() {
	recordService := record.NewService(a.infra.Storage)

	rh := handler.NewRecordHandler(recordService)

	otelShutdown, err := observability.SetupOtelSDK(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup observability: %v", err)
	}

	a.deps = &dependencies{
		RecordHandler: rh,
		OtelShutdown:  otelShutdown,
	}
}

func (a *App) setupServer() {
	a.Server = gin.New()
	a.Server.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	rootGroup := a.Server.Group("/")
	rootGroup.Static("/internal/public", "./internal/public")
	rootGroup.GET("/", func(ctx *gin.Context) {
		ctx.File("./internal/public/index.html")
	})

	baseGroup := a.Server.Group("/api")
	setupHealthCheckRoute(baseGroup)

	a.setupRoutes(baseGroup)
}

func (a *App) startServer() {
	if err := a.Server.Run(fmt.Sprintf(":%s", getPortFallback("PORT", "3000"))); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (a *App) stopApp() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.deps.OtelShutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown observability: %v", err)
	}
}

func (a *App) StartApp() {
	a.startServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	a.stopApp()
}

func setupHealthCheckRoute(g *gin.RouterGroup) {
	g.GET("/health-check", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
}

func getPortFallback(env string, fallback string) string {
	port := os.Getenv(env)
	if port == "" {
		return fallback
	}
	return port
}
