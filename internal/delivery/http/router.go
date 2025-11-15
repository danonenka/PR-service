package http

import (
	"pr-service-task/internal/delivery/http/handlers"
	"pr-service-task/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler       *handlers.UserHandler
	teamHandler       *handlers.TeamHandler
	prHandler         *handlers.PRHandler
	statisticsHandler *handlers.StatisticsHandler
}

func NewRouter(
	userUsecase *usecase.UserUsecase,
	teamUsecase *usecase.TeamUsecase,
	prUsecase *usecase.PRUsecase,
	statisticsUsecase *usecase.StatisticsUsecase,
) *Router {
	return &Router{
		userHandler:       handlers.NewUserHandler(userUsecase, prUsecase, teamUsecase),
		teamHandler:       handlers.NewTeamHandler(teamUsecase),
		prHandler:         handlers.NewPRHandler(prUsecase),
		statisticsHandler: handlers.NewStatisticsHandler(statisticsUsecase),
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	engine.POST("/team/add", r.teamHandler.AddTeam)
	engine.GET("/team/get", r.teamHandler.GetTeam)

	engine.POST("/users/setIsActive", r.userHandler.SetIsActive)
	engine.GET("/users/getReview", r.userHandler.GetReview)

	engine.POST("/pullRequest/create", r.prHandler.CreatePR)
	engine.POST("/pullRequest/merge", r.prHandler.MergePR)
	engine.POST("/pullRequest/reassign", r.prHandler.ReassignReviewer)
}
