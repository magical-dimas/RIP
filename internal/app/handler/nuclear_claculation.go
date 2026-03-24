package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetCalc(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)
	isDraft, err := h.Repository.IsDraftCalc(id, creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	items, calc, err := h.Repository.GetCalc(id, creatorID)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"items":               items,
		"nuclear_calculation": calc,
		"calc_id":             id,
		"minioUrl":            h.Config.MinioURL,
	})
}

func (h *Handler) AddToCalc(ctx *gin.Context) {
	modelIDStr := ctx.PostForm("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)

	err = h.Repository.AddModel(uint(modelID), creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}

func (h *Handler) DeleteCalc(ctx *gin.Context) {
	calcIDStr := ctx.PostForm("calc_id")
	calcID, err := strconv.Atoi(calcIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteCalc(uint(calcID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
