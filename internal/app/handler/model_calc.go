package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"rip_project/internal/app/repository"
	"rip_project/internal/app/serializer"
)

func (h *Handler) AddToCalcAPI(ctx *gin.Context) {
	modelIDStr := ctx.Param("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.Repository.GetCreatorID())
	calc, created, err := h.Repository.GetCalcDraft(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	err = h.Repository.AddModel(uint(modelID), creatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(calc)
	completedCount, _ := h.Repository.GetCompletedItemCount(calc.CalcID)
	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/api/nuclear_calculations/%d", calc.CalcID))
		status = http.StatusCreated
	}
	ctx.JSON(status, serializer.CalcToJSON(calc, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) DeleteFromCalc(ctx *gin.Context) {
	modelID, err := strconv.Atoi(ctx.Param("model_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calcID, err := strconv.Atoi(ctx.Param("calc_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.DeleteModelFromCalc(calcID, modelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(calc)
	completedCount, _ := h.Repository.GetCompletedItemCount(calc.CalcID)
	ctx.JSON(http.StatusOK, serializer.CalcToJSON(calc, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) EditInCalc(ctx *gin.Context) {
	modelID, err := strconv.Atoi(ctx.Param("model_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calcID, err := strconv.Atoi(ctx.Param("calc_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.ModelCalcJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	item, err := h.Repository.EditModelInCalc(calcID, modelID, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.ModelCalcToJSON(item))
}
