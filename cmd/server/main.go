package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/alidonyaie/decision-llm-engine/internal/api"
	"github.com/alidonyaie/decision-llm-engine/internal/config"
	"github.com/alidonyaie/decision-llm-engine/internal/engine"
	"github.com/alidonyaie/decision-llm-engine/internal/llm"
)

const requestTimeout = 20 * time.Second

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not loaded: %v", err)
	}

	appConfig := config.LoadFromEnv()

	promptBuilder, err := engine.NewPromptBuilderFromFile(appConfig.Server.PromptPath)
	if err != nil {
		log.Fatalf("load prompt template: %v", err)
	}

	decisionEngine := engine.NewDecisionEngine(promptBuilder, llm.NewClientFromConfig(appConfig.LLM))
	handler := api.NewHandler(decisionEngine)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(timeoutMiddleware(requestTimeout))
	handler.Register(router)

	server := &http.Server{
		Addr:              ":" + appConfig.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("decision server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("start server: %v", err)
	}
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestContext, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(requestContext)
		c.Next()
	}
}
