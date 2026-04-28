package handler

import (
	"net/http"
	"strconv"

	"rip_project/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetModels(ctx *gin.Context) {
	var models []repository.Model
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		models, err = h.Repository.GetModels()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		models, err = h.Repository.GetModelByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	//idReq := ctx.Param("id")
	//id, err := strconv.Atoi(idReq)
	//if err != nil {
	//	logrus.Error(err)
	//}

	calc, err := h.Repository.GetCalc(1)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"models": models,
		"query":  searchQuery,
		"count":  len(calc.Models),
	})
}

func (h *Handler) GetModel(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	model, err := h.Repository.GetModel(id)
	if err != nil {
		logrus.Error(err)
	}

	reqMod, err := h.Repository.GetCalcForModel(id)

	ctx.HTML(http.StatusOK, "details.html", gin.H{
		"model":  model,
		"reqMod": reqMod,
	})
}

func (h *Handler) GetCalc(ctx *gin.Context) {
	calc, err := h.Repository.GetCalc(1)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"nuclear_calculation": calc,
	})
}
