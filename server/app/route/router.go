package route

import (
     "github.com/gin-gonic/gin"
     "app/modules"
     "net/http"
     "app/sessionInfo"
     "github.com/gin-contrib/sessions"
     "github.com/gin-contrib/sessions/cookie"
	 "log"
    //  "golang.org/x/crypto/bcrypt"
	 "fmt"
     "strconv"
     "app/constants"
     // "app/router"
 )


 var Pagename ="物品一覧"
//  セッション管理に使う

 var LoginInfo sessioninfo.SessionInfo
 var equipments []constants.Equipment
 var Equips_info  []constants.Equip_info
 var students []constants.Student
 var classifications []constants.Classification
 var current_student constants.Student

func InitializeRoutes(router *gin.Engine){
    router.Static("/static", "./static")
    router.LoadHTMLGlob("templates/*.html")

    // セッションの設定
    store := cookie.NewStore([]byte("secret"))
    router.Use(sessions.Sessions("ISDL_IOU", store))

    //登録画面
    router.GET("/signup", func(ctx *gin.Context){
        ctx.HTML(200, "signup.html","dummy")
    })

    // ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        if err :=  modules.Register_user(c); err != nil {
            fmt.Print(err)
        }
        c.Redirect(302, "signup")
    })

    //ログイン画面
    router.GET("/", func(ctx *gin.Context){
        Islogin:=true
        ctx.HTML(200, "login.html", gin.H{"Islogin":Islogin})
    })

    // ユーザーログイン
    router.POST("/login", func(ctx *gin.Context) {
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
            fmt.Printf("%T\n", session) // int
            session.Set("Student_id",student_id )
            session.Set("Name",student_name)
            session.Set("Is_superuser",student_Issuperuser)
            session.Save()

            Equips_info=modules.Update_EquipInfo()  //物品情報の更新

            ctx.Redirect(http.StatusMovedPermanently, "/menu/gallery")
        }   
    })
    
    //ログアウト画面
    router.GET("/logout",func(c *gin.Context){
        //セッションからデータを破棄する
        session := sessions.Default(c)
        session.Clear()
        log.Println("クリア処理")
        session.Save()
        c.Redirect(302, "/")
    })

    //セッション後のページはここにかく
    menu := router.Group("/menu")

    menu.Use(sessionCheck())
    {
        //物品一覧画面
        menu.GET("/gallery", func(ctx *gin.Context){
            root_links:=[]constants.Root_link{}
            page_info :=constants.Page_info{"物品一覧","/menu/gallery",root_links,"Equipments",current_student.Is_superuser}
            ctx.HTML(http.StatusOK, "gallery.html", gin.H{"page_info":page_info,"details":Equips_info})
        })   
    
        //所有物画面
        menu.GET("/mylist", func(ctx *gin.Context){
            fmt.Print(Pagename)
            root_links:=[]constants.Root_link{}
            page_info :=constants.Page_info{"Mylist","/menu/mylist",root_links,"My List",current_student.Is_superuser}
            username:=modules.Get_studentinfo_bySession(ctx,"Name").(string)
            ctx.HTML(200, "mylist.html", gin.H{"page_info":page_info,"username":username,"details":Equips_info})
        })

        //パスワード変更画面
        menu.GET("/changepass", func(ctx *gin.Context){
            fmt.Print(Pagename)
            root_links:=[]constants.Root_link{}
            page_info :=constants.Page_info{"パスワード変更","/menu/changepass",root_links,"Change Password",current_student.Is_superuser}
            ctx.HTML(200, "changepass.html", gin.H{"page_info":page_info,"isGet":true,"isEqual_originpass":true,"isEqual_newpass":true})
        })
        
        menu.POST("/changepass", func(ctx *gin.Context){
            session:=sessions.Default(ctx)
            current_student_id:=session.Get("Is_superuser").(string)
            original_pass:=ctx.PostForm("Original_pass")
            new_pass1:=ctx.PostForm("New_pass1")
            new_pass2:=ctx.PostForm("New_pass2")
            isEqual_originpass,isEqual_newpass:=modules.Change_passwords(current_student_id, original_pass,new_pass1,new_pass2)
            root_links:=[]constants.Root_link{}
            page_info :=constants.Page_info{"パスワード変更","/menu/changepass",root_links,"Change Password",current_student.Is_superuser}
            ctx.HTML(200, "changepass.html", gin.H{"page_info":page_info,"isGet":false,"isEqual_originpass":isEqual_originpass,"isEqual_newpass":isEqual_newpass})
        })

        menu.POST("/rent_apply", func(ctx *gin.Context){
            rental_equipment := modules.Get_RentalApplyEquipment(ctx)
            modules.Rent_apply(rental_equipment)
            // go modules.Send_mail(studentname,equipment_name)
            Equips_info=modules.Update_EquipInfo()  //物品情報の更新
            ctx.Redirect(302, "/menu/gallery")
        })
    
        menu.POST("/return_apply", func(ctx *gin.Context){
            return_equipment := modules.Get_ReturnApplyEquipment(ctx)
            modules.Return_apply(return_equipment)  
            Equips_info=modules.Update_EquipInfo()  //物品情報の更新
            ctx.Redirect(302, "/menu/mylist")
        })  

	}

    sec := router.Group("/sec")
    sec.Use(secCheck())
    {
        //管理画面
		sec.GET("/uadmin", func(ctx *gin.Context) {
			fmt.Print(Pagename)
            root_links:=[]constants.Root_link{}
            page_info :=constants.Page_info{"Admin","/sec/uadmin",root_links,"Admin",current_student.Is_superuser}
			students = modules.Get_userinfo()
			classifications = modules.Get_classinfo()
			ctx.HTML(200, "admin.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info})
		})

        sec.POST("/rent_accept", func(ctx *gin.Context) {
            fmt.Print(Pagename)
            rent_equipment:=modules.Get_RentalAcceptEquipment(ctx)
            modules.Accept_request(rent_equipment)
            Equips_info=modules.Update_EquipInfo()  //物品情報の更新
            ctx.Redirect(302, "/sec/uadmin")
        })

        //ユーザリスト画面
        sec.GET("/list_user", func(c *gin.Context) {
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"list_user","/sec/list_user",root_links,"User List",current_student.Is_superuser}
            students = modules.Get_userinfo()
            c.HTML(200, "list_user.html", gin.H{"page_info":page_info,  "students": students})
        })
    
        sec.POST("/edit_user", func(ctx *gin.Context) {
            fmt.Print(Pagename)
            student := modules.Get_EditUser(ctx)
            delete :=ctx.PostForm("Isdelete")
            isdelete,_:=strconv.ParseBool(delete)
            change :=ctx.PostForm("Ischange")
            ischange,_:=strconv.ParseBool(change)
            if isdelete {
                modules.Delete_userinfo(student)
            }
            if ischange {
                modules.Update_userinfo_admin(student)
            }
            ctx.Redirect(302, "/sec/list_user")
        })
    
        //カテゴリーリスト画面
        sec.GET("/list_classification", func(c *gin.Context) {
            fmt.Print(Pagename)
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"list_classification","/sec/list_classification",root_links,"Class List",current_student.Is_superuser}
            c.HTML(200, "list_classification.html", gin.H{"page_info":page_info, "classifications": classifications})
        })	
    
        sec.POST("/edit_classification", func(c *gin.Context) {
            fmt.Print(Pagename)
            id,_ := strconv.Atoi(c.PostForm("Id"))
            name := c.PostForm("Name")
            classification := constants.Classification{id,name}
            delete :=c.PostForm("Isdelete")
            isdelete,_:=strconv.ParseBool(delete)
            change :=c.PostForm("Ischange")
            ischange,_:=strconv.ParseBool(change)
            if isdelete {
                modules.Delete_classinfo(classification) 
            }
            if ischange {
                 modules.Update_classinfo_admin(classification)
            }
            classifications = modules.Get_classinfo()
            c.Redirect(302, "/sec/list_classification")
        })

        //ユーザ追加画面
        sec.GET("/add_user", func(ctx *gin.Context){
            isexists_userform:=constants.Isexists_userform{true,true,true,true}
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_user","/sec/add_user",root_links,"Add User",current_student.Is_superuser}
            ctx.HTML(200, "add_user.html", gin.H{"page_info":page_info,  "students": students,"isexists_userform":isexists_userform})
        })
    
        sec.POST("/add_user", func(c *gin.Context) {
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_user","/sec/add_user",root_links,"Add User",current_student.Is_superuser}
            isexists_userform,isfull_form :=modules.Isexists_userform(c)
            if isfull_form{
            if err :=  modules.Register_user(c); err != nil {
                fmt.Print(err)
            }
        } 
        c.HTML(200, "add_user.html", gin.H{"page_info":page_info,  "students": students, "isfull_form":isfull_form,"isexists_userform": isexists_userform})
        })

        // 物品追加画面
        sec.GET("/add_equipment", func(c *gin.Context) {
            fmt.Print(Pagename)
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_equipment","/sec/add_equipment",root_links,"Add Equipment",current_student.Is_superuser}
            isexists_equipform:=constants.Isexists_equipform{true,true}

            c.HTML(200, "add_equipment.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info, "classifications": classifications,"isexists_equipform": isexists_equipform})

        })
        sec.POST("/add_equipment", func(c *gin.Context) {
            fmt.Print(Pagename)
          
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_equipment","/sec/add_equipment",root_links,"Add Equipment",current_student.Is_superuser}
            name := c.PostForm("Name")
            classification_id,_ := strconv.Atoi(c.PostForm("Classifications_id"))
 
            isexists_equipform,isfull_form :=modules.Isexists_equipform(c)
            if isfull_form{
            modules.Insert_equipinfo(name,classification_id)
            Equips_info=modules.Update_EquipInfo()  //物品情報の更新
            }
            c.HTML(200, "add_equipment.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info, "classifications": classifications,"isfull_form":isfull_form,"isexists_equipform": isexists_equipform})

        })

        // ジャンル追加画面
        sec.GET("/add_classification", func(c *gin.Context) {
            fmt.Print(Pagename)
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_classification","/sec/add_classification",root_links,"Add Class",current_student.Is_superuser}

            c.HTML(200, "add_classification.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info, "classifications": classifications,"isexists_classificationsform":true})

        })
        sec.POST("/add_classification", func(c *gin.Context) {
            fmt.Print(Pagename)
            name := c.PostForm("Name")
            isexists_classificationsform,isfull_form:=modules.Isexists_classificationsform(c)
            if isfull_form {
               modules.Insert_classificationinfo(name) 
            }
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"add_classification","/sec/add_classification",root_links,"Add Class",current_student.Is_superuser}

            c.HTML(200, "add_classification.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info, "classifications": classifications,"isfull_form":isfull_form,"isexists_classificationsform":isexists_classificationsform})
        })

        // 物品編集画面
        sec.GET("/list_equipment", func(c *gin.Context) {
            fmt.Print(Pagename)	
            root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
            page_info :=constants.Page_info{"list_equipment","/sec/list_equipment",root_links,"Equipment List",current_student.Is_superuser}
            c.HTML(200, "list_equipment.html", gin.H{"page_info":page_info,"equipments": equipments,"details":Equips_info, "classifications": classifications})
        })
    
        sec.POST("/edit_equipment", func(ctx *gin.Context) {
            fmt.Print(Pagename)
            equipment := modules.Get_EditEquipment(ctx)
            delete :=ctx.PostForm("Isdelete")
            isdelete,_:=strconv.ParseBool(delete)
            change :=ctx.PostForm("Ischange")
            ischange,_:=strconv.ParseBool(change)
            if isdelete {
               modules.Delete_equipinfo(equipment) 
            }  
            if ischange {
                modules.Update_equipinfo(equipment) 
             }
            Equips_info=modules.Update_EquipInfo()  //物品情報の更新
            ctx.Redirect(302, "/sec/list_equipment")
        })
    }
}



func sessionCheck() gin.HandlerFunc {
    return func(c *gin.Context) {

        session := sessions.Default(c)
        LoginInfo.Student_id = session.Get("Student_id")
        LoginInfo.Name = session.Get("Name")
        LoginInfo.Is_superuser = session.Get("Is_superuser")

        // セッションがない場合、ログインフォームをだす
        if LoginInfo.Student_id == nil {
            log.Println("ログインしていません")
            c.Redirect(http.StatusMovedPermanently, "/login")
            c.Abort() // これがないと続けて処理されてしまう
        } else {
            c.Set("Student_id", LoginInfo.Student_id.(string))
            c.Set("Name", LoginInfo.Name.(string))
            c.Set("Is_superuser", LoginInfo.Is_superuser.(string))

            c.Next()
            
            
        }
        log.Println("ログインチェック終わり")
    }


}

func secCheck() gin.HandlerFunc {

    return func(c *gin.Context) {
        session := sessions.Default(c)
        LoginInfo.Student_id = session.Get("Student_id")
        LoginInfo.Name = session.Get("Name")
        LoginInfo.Is_superuser = session.Get("Is_superuser")


        // セッションがない場合、ログインフォームをだす
        if LoginInfo.Is_superuser != true {
            log.Println("uadminアカウントではありません")
            // c.Redirect(http.StatusMovedPermanently, "/menu/gallery")
            c.Abort() // これがないと続けて処理されてしまう
        } else {
     
            c.Set("Student_id", LoginInfo.Student_id.(string))
            c.Set("Name", LoginInfo.Name.(string))
            c.Set("Is_superuser", LoginInfo.Is_superuser.(string))
  
            c.Next()
            
        }

    }


}