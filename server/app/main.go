package main

import _"fmt"
import (
   _"net/http"
   _"github.com/labstack/echo/v4"
   _"github.com/labstack/echo/v4/middleware"
	_"app/modules"
	"github.com/gin-gonic/gin"
	"app/route"
	"github.com/thinkerou/favicon"
)

var router *gin.Engine


func main(){
	//routing
	router := gin.Default()
	router.Use(favicon.New("./favicon.ico")) 
	//router.goを反映させる
	route.InitializeRoutes(router)
    router.Run(":5050")


}	


