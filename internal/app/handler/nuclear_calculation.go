package handler

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"rip_project/internal/app/repository"
	"rip_project/internal/app/serializer"
)

// @Summary Получение иконки корзины (сводка черновика)
// @Tags Заявки
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/nuclear_calculations/items [get]
func (h *Handler) GetCalcItems(ctx *gin.Context) {
	creatorID := uint(h.getEngineerID(ctx))
	count := h.Repository.GetCalcModelCount(creatorID)
	if count == 0 {
		calc, err := h.Repository.CheckCurrentDraft(creatorID)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":      "no_draft",
				"model_count": 0,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":          calc.CalcID,
			"model_count": 0,
		})
		return
	}
	calcID := h.Repository.GetActiveCalcID(creatorID)
	ctx.JSON(http.StatusOK, gin.H{
		"id":          calcID,
		"model_count": count,
	})
}

// @Summary Список всех сформированных заявок (с фильтрами)
// @Tags Заявки
// @Security BearerAuth
// @Produce json
// @Param status query string false "Фильтр по статусу" Enums(draft, pending, approved, rejected)
// @Param from-date query string false "Дата от (формат: YYYY-MM-DD)" format(date) example("2024-01-01")
// @Param to-date query string false "Дата до (формат: YYYY-MM-DD)" format(date) example("2024-12-31")
// @Success 200 {array} map[string]interface{} "Список расчётов"
// @Router /api/nuclear_calculations [get]
func (h *Handler) GetAllCalcs(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from, to time.Time
	if fromDate != "" {
		t, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = t
	}
	toDate := ctx.Query("to-date")
	if toDate != "" {
		t, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = t
	}
	status := ctx.Query("status")
	calcs, err := h.Repository.GetAllCalcs(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.CalcJSON, 0, len(calcs))
	for _, calc := range calcs {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(calc)
		completedCount, _ := h.Repository.GetCompletedItemCount(calc.CalcID)
		resp = append(resp, serializer.CalcToJSON(calc, creatorLogin, moderatorLogin, completedCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Получение заявки по ID
// @Tags Заявки
// @Security BearerAuth
// @Param id path int true "ID Заявки"
// @Success 200 {object} ds.Nuclear_calculation
// @Router /api/nuclear_calculations/{id} [get]
func (h *Handler) GetCalcAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.GetSingleCalc(id)
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

	items, err := h.Repository.GetCalcItems(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(calc)
	completedCount, _ := h.Repository.GetCompletedItemCount(calc.CalcID)
	itemsResp := make([]serializer.ModelCalcJSON, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, serializer.ModelCalcToJSON(item))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"calc":   serializer.CalcToJSON(calc, creatorLogin, moderatorLogin, completedCount),
		"models": itemsResp,
	})
}

// @Summary Изменение заявки
// @Tags Заявки
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path integer true "ID заявки" minimum(1)
// @Param input body serializer.CalcJSON true "Обновляемые поля заявки"
// @Success 200 "Успешно"
// @Router /api/nuclear_calculations/{id} [put]
func (h *Handler) EditCalc(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.CalcJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.EditCalc(id, j)
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

// @Summary Сформировать заявку (запуск расчетов)
// @Description Переводит заявку из draft в сформирован
// @Tags Заявки
// @Security BearerAuth
// @Param id path int true "ID Заявки"
// @Success 200 {object} map[string]interface{}
// @Router /api/nuclear_calculations/{id}/form [put]
func (h *Handler) FormCalc(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.GetSingleCalc(id)
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
	if calc.CreatorID != h.getEngineerID(ctx) {
		h.errorHandler(ctx, http.StatusForbidden, repository.ErrNotAllowed)
	}
	calc, err = h.Repository.FormCalc(id)
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
	completedCount, _ := h.Repository.GetCompletedItemCount(calc.CalcID)
	items, err := h.Repository.GetCalcItems(int(calc.CalcID))
	itemsResp := make([]serializer.ModelCalcJSON, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, serializer.ModelCalcToJSON(item))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"calculated": completedCount,
		"items": itemsResp,
	})
}

// @Summary Завершение или отклонение заявки
// @Description Доступно только Технику
// @Tags Заявки
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path integer true "ID заявки" minimum(1)
// @Param input body serializer.StatusJSON true "Действие: завершить или отклонить"
// @Success 200 {object} map[string]string
// @Router /api/nuclear_calculations/{id}/finish [put]
func (h *Handler) FinishCalc(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	calc, err := h.Repository.FinishCalc(id, statusJSON.Status, int(h.getEngineerID(ctx)))
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

// @Summary Логическое удаление заявки
// @Tags Заявки
// @Security BearerAuth
// @Param id path integer true "ID Заявки"
// @Success 200 "Успешно"
// @Router /api/nuclear_calculations/{id} [delete]
func (h *Handler) DeleteCalcAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calc, err := h.Repository.GetSingleCalc(id)
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
	if calc.CreatorID != h.getEngineerID(ctx) {
		h.errorHandler(ctx, http.StatusForbidden, repository.ErrNotAllowed)
	}
	_, err = h.Repository.DeleteCalc(id)
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}

func (h *Handler) GetCalc(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(h.getEngineerID(ctx))
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

func (h *Handler) DeleteCalc(ctx *gin.Context) {
	calcIDStr := ctx.PostForm("calc_id")
	calcID, err := strconv.Atoi(calcIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	_, err = h.Repository.DeleteCalc(calcID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) AddToCalc(ctx *gin.Context) {
	modelIDStr := ctx.PostForm("model_id")
	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	creatorID := uint(h.getEngineerID(ctx))

	err = h.Repository.AddModel(uint(modelID), creatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, ctx.Request.Referer())
}
