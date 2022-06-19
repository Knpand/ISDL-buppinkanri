package handler

import (
	"github.com/gin-gonic/gin"
	"app/modules"
	"net/http"
	"github.com/gin-contrib/sessions"
	"log"
)


func DisplayLoginPage(ctx *gin.Context){
	Islogin:=true
	ctx.HTML(200, "login.html", gin.H{"Islogin":Islogin})
}

func DisplaySignupPage(ctx *gin.Context){
	ctx.HTML(200, "signup.html","dummy")
}

func RegisterUser(c *gin.Context) {
	if err :=  modules.Register_user(c); err != nil {
		log.Print(err)
	}
	c.Redirect(302, "signup")
}

func Login(ctx *gin.Context) {
	student_id:=ctx.PostForm("student_id")
	usr_FormPassword:=ctx.PostForm("password")
	usr_DBPassword:=modules.GetDBUserPassword(student_id)
	if err := modules.CompareHashAndPassword(usr_DBPassword, usr_FormPassword); err != nil {
		log.Println("err occured")
		log.Println(err)
		Islogin:=false
		ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"Islogin":Islogin})
		ctx.Abort()
	} else {
		log.Println("log in")
		student_info := modules.Fetch_studentinfo_byID(student_id)
		student_name := student_info.Name
		student_Issuperuser := student_info.Is_superuser
		session := sessions.Default(ctx)    //ユーザに紐づくセッションを取得
	
		session.Set("Student_id",student_id )
		session.Set("Name",student_name)
		session.Set("Is_superuser",student_Issuperuser)
		log.Println(session.Get("Student_id" ))
		session.Save()

		ctx.Redirect(http.StatusMovedPermanently, "/menu/gallery")
	}   
}

func Logout(c *gin.Context){
	//セッションからデータを破棄する
	session := sessions.Default(c)
	session.Clear()
	log.Println("クリア処理")
	session.Save()
	c.Redirect(302, "/")
}

func SessionCheck() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        Student_id := session.Get("Student_id")
        // セッションがない場合、ログインフォームをだす
        if Student_id == nil {
            log.Println("ログインしていません")
            ctx.Redirect(http.StatusMovedPermanently, "/")
            ctx.Abort() // これがないと続けて処理されてしまう
        } 
        log.Println("ログインチェック終わり")
    }
}


func SecCheck() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        session := sessions.Default(ctx)
        Is_superuser := session.Get("Is_superuser")
        // セッションがない場合、ログインフォームをだす
        if Is_superuser != true {
            log.Println("uadminアカウントではありません")
            ctx.Redirect(http.StatusMovedPermanently, "/")
            ctx.Abort() // これがないと続けて処理されてしまう
        } 
    }

}