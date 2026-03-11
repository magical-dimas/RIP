package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(1)
	isDraft, err := h.Repository.IsDraftRequest(id, creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !isDraft {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	items, req, err := h.Repository.GetRequest(id, creatorID)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"items":      items,
		"request":    req,
		"request_id": id,
		"minioUrl":   h.Config.MinioURL,
	})
}

func (h *Handler) AddToRequest(ctx *gin.Context) {
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

func (h *Handler) DeleteRequest(ctx *gin.Context) {
	reqIDStr := ctx.PostForm("request_id")
	reqID, err := strconv.Atoi(reqIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteRequest(uint(reqID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
