package modules

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/joho/godotenv"
	"os"
	"log"
)
func SessionSetting(router *gin.Engine){
    // セッションの設定
	godotenv.Load(".env")
	log.Print("SessionSetting")
	log.Print(os.Getenv("IOU_key"))
    secret := os.Getenv("IOU_key")
    store := cookie.NewStore([]byte(secret))
    router.Use(sessions.Sessions("ISDL_IOU", store))

}


