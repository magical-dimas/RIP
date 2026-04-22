package handler

import (
	"errors"

	"rip_project/internal/app/config"
	"rip_project/internal/app/repository"
	"rip_project/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
}

func NewHandler(r *repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	//website
	router.GET("/", h.GetModels)
	router.GET("/model/:id", h.GetModel)
	router.GET("/nuclear_calculations/:id", h.GetCalc)
	router.POST("/nuclear_calculations/add", h.AddToCalc)
	router.POST("/nuclear_calculations/delete", h.DeleteCalc)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")

	//unauthorized
	models := api.Group("/models")
	{
		models.GET("", h.GetModelsAPI)
		models.GET("/:id", h.GetModelAPI)
	}
	engineers := api.Group("/engineers")
	{
		engineers.POST("/register", h.CreateEngineer)
		engineers.POST("/login", h.SignIn)
	}

	optionalAuth := api.Group("/")
	optionalAuth.Use(h.OptionalAuthCheck(role.Engineer))
	{
		optionalAuth.GET("/nuclear_calculations/items", h.GetCalcItems)
	}

	//authorized
	authGroup := api.Group("/")
	authGroup.Use(h.WithAuthCheck(role.Engineer))
	{
		authGroup.POST("/engineers/logout", h.SignOut)

		authGroup.POST("/model_calculation/add/:model_id", h.AddToCalcAPI)
		authGroup.PUT("/model_calculation/:model_id", h.EditInCalc)
		authGroup.DELETE("/model_calculation/:model_id", h.DeleteFromCalc)

		authGroup.GET("/nuclear_calculations", h.GetAllCalcs)
		authGroup.GET("/nuclear_calculations/:id", h.GetCalcAPI)
		authGroup.PUT("/nuclear_calculations/:id", h.EditCalc)
		authGroup.PUT("/nuclear_calculations/:id/form", h.FormCalc)
		authGroup.DELETE("/nuclear_calculations/:id", h.DeleteCalcAPI)
	}

	//moderator
	modGroup := api.Group("/")
	modGroup.Use(h.WithAuthCheck(role.Technician))
	{
		modGroup.POST("/models", h.CreateModel)
		modGroup.PUT("/nuclear_calculations/:id/finish", h.FinishCalc)
	}
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}
