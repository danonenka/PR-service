package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	httphandler "pr-service-task/internal/delivery/http"
	"pr-service-task/internal/repository/postgres"
	"pr-service-task/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "pr_service")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)
	assignmentRepo := postgres.NewReviewerAssignmentRepository(db)

	reassignmentUsecase := usecase.NewReassignmentUsecase(prRepo, userRepo, assignmentRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, teamRepo, reassignmentUsecase)
	teamUsecase := usecase.NewTeamUsecase(teamRepo, userRepo)
	prUsecase := usecase.NewPRUsecase(prRepo, userRepo, assignmentRepo)
	statisticsUsecase := usecase.NewStatisticsUsecase(prRepo, assignmentRepo, userRepo)

	router := httphandler.NewRouter(userUsecase, teamUsecase, prUsecase, statisticsUsecase)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	// Swagger UI
	openapiPath := "/app/openapi.yaml"
	if _, err := os.Stat("./openapi.yaml"); err == nil {
		openapiPath = "./openapi.yaml"
	}
	engine.StaticFile("/swagger/openapi.yaml", openapiPath)
	
	// Swagger UI HTML
	engine.GET("/swagger-ui", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html>
<head>
	<title>API Documentation</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
	<script>
		window.onload = function() {
			SwaggerUIBundle({
				url: "/swagger/openapi.yaml",
				dom_id: '#swagger-ui',
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIBundle.presets.standalone
				]
			});
		};
	</script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})
	
	router.SetupRoutes(engine)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on 0.0.0.0:%s", port)
	log.Printf("Swagger UI available at http://localhost:%s/swagger-ui", port)
	log.Printf("OpenAPI spec available at http://localhost:%s/swagger/openapi.yaml", port)
	if err := engine.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
