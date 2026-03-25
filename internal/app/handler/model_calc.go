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

// @Summary Добавление модели в корзину (черновик)
// @Tags Корзина
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param model_id path integer true "ID модели" minimum(1)
// @Success 201 {object} map[string]string
// @Router /api/model_calculation/add/{model_id} [post]
func (h *Handler) AddToCalcAPI(ctx *gin.Context) {
	modelIDStr := ctx.Param("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.getEngineerID(ctx))
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

// @Summary Удаление элемента из корзины
// @Tags Корзина
// @Security BearerAuth
// @Produce json
// @Param model_id path integer true "ID модели" minimum(1)
// @Success 200 "Успешно удалено"
// @Router /api/model_calculation/{model_id}/{calc_id} [delete]
func (h *Handler) DeleteFromCalc(ctx *gin.Context) {
	modelID, err := strconv.Atoi(ctx.Param("model_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.getEngineerID(ctx))
	calc, _, err := h.Repository.GetCalcDraft(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	calc, err = h.Repository.DeleteModelFromCalc(int(calc.CalcID), modelID)
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

// @Summary Обновление элемента в корзине
// @Tags Корзина
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param model_id path integer true "ID модели" minimum(1)
// @Param input body serializer.ModelCalcJSON true "Новые параметры"
// @Success 200 "Успешно"
// @Router /api/model_calculation/{model_id}/{calc_id} [put]
func (h *Handler) EditInCalc(ctx *gin.Context) {
	modelID, err := strconv.Atoi(ctx.Param("model_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.getEngineerID(ctx))
	calc, _, err := h.Repository.GetCalcDraft(creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	var j serializer.ModelCalcJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	item, err := h.Repository.EditModelInCalc(int(calc.CalcID), modelID, j)
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
