package main

import (
	"blog/config"
	"blog/models"
	"blog/router"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

func init() {
	config.Init()
	models.Init()
}

func main() {

	gin.SetMode("debug")

	builder := strings.Builder{}
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(config.App.Server.Port))

	r := router.Init()
	err := r.Run(builder.String())

	if err != nil {
		log.Fatalf("main.ListenAndServe: %v", err)
	}
}
