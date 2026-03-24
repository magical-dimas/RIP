package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rip_project/internal/app/config"
	"rip_project/internal/app/dsn"
	"rip_project/internal/app/handler"
	"rip_project/internal/app/repository"
	"rip_project/internal/pkg"
)

func main() {
	router := gin.Default()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	logrus.Info("DSN: ", postgresString)

	rep, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	hand := handler.NewHandler(rep, conf)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
