package handler

import (
	"fmt"
	"lab1_isbn/internal/core/domain"
	"lab1_isbn/internal/core/ports/services"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const name = "Record"

type (
	pathParams struct {
		Number uint `uri:"number" binding:"required"`
	}

	bodyParams struct {
		Name        string `form:"name" binding:"required"`
		Description string `form:"description"`
	}

	RecordHandler struct {
		RecordService services.RecordService
		observability
	}

	observability struct {
		tracer        trace.Tracer
		meter         metric.Meter
		logger        slog.Logger
		recordCounter metric.Int64Counter
	}
)

func NewRecordHandler(rs services.RecordService) *RecordHandler {
	var err error

	observability := observability{
		tracer: otel.Tracer(name),
		meter:  otel.Meter(name),
		logger: *otelslog.NewLogger(name),
	}

	observability.recordCounter, err = observability.
		meter.
		Int64Counter(
			"record.counter",
			metric.WithDescription("Number of records created"),
			metric.WithUnit("{number}"),
		)
	if err != nil {
		observability.logger.Error(fmt.Sprintf("Failed to create counter: %s", err.Error()))
	}

	return &RecordHandler{
		RecordService: rs,
		observability: observability,
	}
}

func (rh *RecordHandler) CreateRecord(c *gin.Context) {
	ctx, span := rh.tracer.Start(c.Request.Context(), "CreateRecord")
	defer span.End()

	body := &bodyParams{}

	if err := c.ShouldBind(body); err != nil {
		span.RecordError(err)
		rh.logger.InfoContext(ctx, "Failed to bind body", "error", err.Error(), "body", body)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	task := domain.Task{
		Name: body.Name,
		Desc: body.Description,
	}

	if err := rh.RecordService.CreateRecord(ctx, task); err != nil {
		span.RecordError(err)
		rh.logger.InfoContext(ctx, "Failed to create record", "error", err.Error(), "body", body)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Create Record",
	})
}

func (rh *RecordHandler) UpdateRecord(c *gin.Context) {
	ctx, span := rh.tracer.Start(c.Request.Context(), "CreateRecord")
	defer span.End()

	body := &bodyParams{}

	if err := c.ShouldBind(body); err != nil {
		span.RecordError(err)
		rh.logger.InfoContext(ctx, "Failed to bind body", "error", err.Error(), "body", body)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	task := domain.Task{
		Name: body.Name,
		Desc: body.Description,
	}

	if err := rh.RecordService.UpdateRecord(ctx, task); err != nil {
		span.RecordError(err)
		rh.logger.InfoContext(ctx, "Failed to update record", "error", err.Error(), "body", body)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update Record",
	})
}
