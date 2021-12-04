package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Ping PingTest = &ping{}
)

type PingTest interface {
	Ping(ctx *gin.Context)
}

type ping struct{}

func (p *ping) Ping(ctx *gin.Context) {
	serverMessage := map[string]interface{}{
		"status":  200,
		"code":    "server_working",
		"message": "server is in good condition",
	}

	ctx.JSON(http.StatusOK, serverMessage)
}
