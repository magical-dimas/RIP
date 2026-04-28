package handler

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"rip_project/internal/app/ds"
	"rip_project/internal/app/repository"
	"rip_project/internal/app/serializer"
)

// @Summary Получение всех услуг
// @Description Возвращает список всех моделей реакторов из каталога
// @Tags Услуги
// @Produce json
// @Param search query string false "Поиск по названию"
// @Success 200 {array} ds.Model
// @Router /api/models [get]
func (h *Handler) GetModelsAPI(ctx *gin.Context) {
	var models []ds.Model
	var err error

	searchQuery := ctx.Query("Title")
	if searchQuery == "" {
		models, err = h.Repository.GetModels()
	} else {
		models, err = h.Repository.GetModelsByTitle(searchQuery)
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.ModelJSON, 0, len(models))
	for _, m := range models {
		resp = append(resp, serializer.ModelToJSON(m))
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Получение услуги по ID
// @Tags Услуги
// @Produce json
// @Param id path int true "ID Модели"
// @Success 200 {object} ds.Model
// @Router /api/models/{id} [get]
func (h *Handler) GetModelAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	model, err := h.Repository.GetModel(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.ModelToJSON(*model))
}

// @Summary Добавление новой услуги
// @Description Доступно только Технику
// @Tags Услуги
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Название"
// @Param description formData string false "Описание"
// @Param short_desc formData string false "Краткое описание"
// @Param power formData number false "Мощность"
// @Param fuel_usage formData number false "Расход Топлива"
// @Success 201 {object} ds.Model
// @Router /api/models [post]
func (h *Handler) CreateModel(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var j serializer.ModelJSON
	if strings.HasPrefix(contentType, "application/json") {
		if err := ctx.BindJSON(&j); err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		title := ctx.PostForm("title")
		desc := ctx.PostForm("description")
		s_desc := ctx.PostForm("short_desc")
		if title == "" {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("title is required"))
			return
		}
		power := float64(55)
		fuel_usage := float64(45)
		if v := ctx.PostForm("power"); v != "" {
			fmt.Sscanf(v, "%f", &power)
		}
		if v := ctx.PostForm("fuel_usage"); v != "" {
			fmt.Sscanf(v, "%f", &fuel_usage)
		}
		j = serializer.ModelJSON{
			Title:       title,
			Description: desc,
			ShortDesc:   s_desc,
			Power:       power,
			FuelUsage:   fuel_usage,
		}
	}

	model, err := h.Repository.CreateModel(j)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if imageFile, err := ctx.FormFile("photo"); err == nil {
		m, err := h.Repository.AddPhoto(ctx, int(model.ModelID), imageFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		model = *m
	}
	if videoFile, err := ctx.FormFile("video"); err == nil {
		m, err := h.Repository.AddVideo(ctx, int(model.ModelID), videoFile)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		model = *m
	}

	ctx.Header("Location", fmt.Sprintf("/api/model/%d", model.ModelID))
	ctx.JSON(http.StatusCreated, serializer.ModelToJSON(model))
}

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
