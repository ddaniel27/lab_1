package api

import "github.com/gin-gonic/gin"

func (a *App) setupRoutes(g *gin.RouterGroup) {
	handler := a.deps.RecordHandler

	recordGroup := g.Group("/tasks")

	recordGroup.POST("", handler.CreateRecord)
	recordGroup.POST("/update", handler.UpdateRecord)
}
