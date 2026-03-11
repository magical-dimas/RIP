package main

import (
	"fmt"
	"strconv"

	"rip_project/internal/app/config"
	"rip_project/internal/app/dsn"
	"rip_project/internal/app/handler"
	"rip_project/internal/app/repository"
	"rip_project/internal/pkg"

	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	router.SetFuncMap(template.FuncMap{
		"printf": fmt.Sprintf,
		"num": func(v interface{}) float64 {
			if v == nil {
				return 0
			}
			switch x := v.(type) {
			case float64:
				return x
			case *float64:
				if x == nil {
					return 0
				}
				return *x
			default:
				f, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
				return f
			}
		},
	})

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
