package handler

import (
	"net/http"
	"rip_project/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetModels(ctx *gin.Context) {
	var models []ds.Model
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		models, err = h.Repository.GetModels()
	} else {
		models, err = h.Repository.GetModelsByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	creatorID := uint(1)
	calcCount := h.Repository.GetCalcModelCount(creatorID)
	activeCalcID := h.Repository.GetActiveCalcID(creatorID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"models":     models,
		"query":      searchQuery,
		"calc_count": calcCount,
		"calc_id":    activeCalcID,
		"minioUrl":   h.Config.MinioURL,
	})
}

func (h *Handler) GetModel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	model, err := h.Repository.GetModel(id)
	if err != nil {
		logrus.Error(err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "details.html", gin.H{
		"model":    model,
		"minioUrl": h.Config.MinioURL,
	})
}
