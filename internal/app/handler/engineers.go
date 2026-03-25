package handler

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"rip_project/internal/app/role"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"rip_project/internal/app/repository"
	"rip_project/internal/app/serializer"
)

// @Summary Регистрация пользователя
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param input body serializer.EngineerJSON true "Данные регистрации"
// @Success 201 {object} ds.Engineers
// @Router /api/engineers/register [post]
func (h *Handler) CreateEngineer(ctx *gin.Context) {
	var j serializer.EngineerJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	engineer, err := h.Repository.CreateEngineer(j)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.Header("Location", fmt.Sprintf("/api/engineers/%d", engineer.EngineerID))
	ctx.JSON(http.StatusCreated, serializer.EngineerToJSON(engineer))
}

// @Summary Авторизация (Login)
// @Description Выдает JWT токен в случае успешной авторизации
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param input body serializer.EngineerJSON true "Данные для входа"
// @Success 200 {object} map[string]interface{}
// @Router /api/engineers/login [post]
func (h *Handler) SignIn(ctx *gin.Context) {
	var j serializer.EngineerJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if j.Login == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'login' is required"))
		return
	}
	if j.Password == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("field 'password' is required"))
		return
	}
	engineer, err := h.Repository.SignIn(j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) || err.Error() == "неверный логин или пароль" {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("invalid login or password"))
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	var rl = role.Guest
	if engineer.IsTechnician {
		rl = role.Technician
	} else {
		rl = role.Engineer
	}
	claims := &JWTClaims{
		EngineerID: engineer.EngineerID,
		Role:       rl,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"expires_in":   expirationTime.Unix(),
		"role":         rl,
	})
}

// @Summary Выход (Logout)
// @Description Помещает переданный JWT в Blacklist Redis'а
// @Tags Пользователи
// @Security BearerAuth
// @Router /api/engineers/logout [post]
func (h *Handler) SignOut(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, jwtPrefix) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Токен не передан"})
		return
	}

	tokenStr := authHeader[len(jwtPrefix):]

	token, _ := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	var expTime time.Duration = 24 * time.Hour
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		expTime = time.Until(time.Unix(claims.ExpiresAt, 0))
	}

	h.Repository.GetRedis().Set(ctx.Request.Context(), "blacklist:"+tokenStr, true, expTime)

	ctx.JSON(http.StatusOK, gin.H{"message": "Вы успешно вышли из системы"})
}
