package api

import (
	"log"
	"rip_project/internal/app/handler"
	"rip_project/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.GetModels)
	r.GET("/model/:id", h.GetModel)
	r.GET("/cart", h.GetCalc)

	r.Run()
	log.Println("Server down")
}
