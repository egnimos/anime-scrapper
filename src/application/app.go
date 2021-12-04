package application

import (
	// "io"
	// "os"

	"fmt"
	"os"

	"github.com/egnimos/anime-scrapper/src/server_engine"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	router *gin.Engine
)

func init() {
	gin.ForceConsoleColor()
	// f, _ := os.Create("network.log")
	// gin.DefaultWriter = io.MultiWriter(f)
	router = gin.Default()
	router.SetTrustedProxies([]string{"192.168.0.1"})
}

func Application() {
	//server app
	loaderr := godotenv.Load("C:/Users/nitee/Documents/workspace/GOLANG-PRACTICE/anime-scrapper/.env")
	if loaderr != nil {
		fmt.Println("error loading .env file")
	}
	
	// C:\Users\nitee\Documents\workspace\GOLANG-PRACTICE\video-scrapper\src\json\servers.json
	serverJsonInfo := &server_engine.ServerJson{
		Path:       os.Getenv("SERVER_JSON"),
		ServerList: server_engine.ServerList{},
	}

	server_engine.Servers.Server(serverJsonInfo)

	//set the router
	routers(router)

	//start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
