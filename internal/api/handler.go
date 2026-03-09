package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/alidonyaie/decision-llm-engine/internal/engine"
	"github.com/alidonyaie/decision-llm-engine/internal/model"
)

// Handler wires HTTP requests to the decision engine.
type Handler struct {
	engine *engine.DecisionEngine
}

// NewHandler creates an API handler.
func NewHandler(decisionEngine *engine.DecisionEngine) *Handler {
	return &Handler{engine: decisionEngine}
}

// Register binds API routes.
func (h *Handler) Register(router gin.IRoutes) {
	router.POST("/v1/decision/analyze", h.handleAnalyze)
	router.GET("/health", h.handleHealth)
	h.registerSwagger(router)
}

func (h *Handler) handleAnalyze(context *gin.Context) {
	var input model.AnalyzeRequest
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(input.Question) == "" {
		context.JSON(http.StatusBadRequest, map[string]string{"error": "question is required"})
		return
	}

	decision, err := h.engine.Analyze(context.Request.Context(), input.Question)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, model.AnalyzeResponse{Decision: decision})
}

func (h *Handler) handleHealth(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
