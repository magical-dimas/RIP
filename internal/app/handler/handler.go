package handler

import (
	"errors"

	"rip_project/internal/app/config"
	"rip_project/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	api := router.Group("/api")

	models := api.Group("/models")
	{
		models.GET("", h.GetModelsAPI)
		models.GET("/:id", h.GetModelAPI)
		models.POST("", h.CreateModel)
	}

	nuc_calcs := api.Group("/nuclear_calculations")
	{
		nuc_calcs.GET("/items", h.GetCalcItems)
		nuc_calcs.GET("", h.GetAllCalcs)
		nuc_calcs.GET("/:id", h.GetCalcAPI)
		nuc_calcs.PUT("/:id", h.EditCalc)
		nuc_calcs.PUT("/:id/form", h.FormCalc)
		nuc_calcs.PUT("/:id/finish", h.FinishCalc)
		nuc_calcs.DELETE("/:id", h.DeleteCalcAPI)
	}

	mc := api.Group("/model_calculation")
	{
		mc.POST("/add/:model_id", h.AddToCalcAPI)
		mc.DELETE("/:model_id", h.DeleteFromCalc)
		mc.PUT("/:model_id", h.EditInCalc)
	}

	users := api.Group("/users")
	{
		users.POST("/register", h.CreateUser)
		users.POST("/login", h.SignIn)
		users.POST("/logout", h.SignOut)
	}

	//website
	router.GET("/", h.GetModels)
	router.GET("/model/:id", h.GetModel)
	router.GET("/nuclear_calculations/:id", h.GetCalc)
	router.POST("/nuclear_calculations/add", h.AddToCalc)
	router.POST("/nuclear_calculations/delete", h.DeleteCalc)
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
