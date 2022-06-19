package route

import (
     "github.com/gin-gonic/gin"
     "app/handler"
     "app/constants"
     "app/modules"
 )

//  セッション管理に使う


 var Equips_info []constants.Equip_info

func InitializeRoutes(router *gin.Engine){
    router.Static("/static", "./static")
    router.LoadHTMLGlob("templates/*.html")

    modules.SessionSetting(router)

    //登録画面
    router.GET("/signup", handler.DisplaySignupPage)
    // ユーザー登録
    router.POST("/signup", handler.RegisterUser)
    //ログイン画面
    router.GET("/", handler.DisplayLoginPage)
    // ログイン処理
    router.POST("/login",handler.Login )
    //ログアウト画面
    router.GET("/logout",handler.Logout)

    //セッション後のページ
    menu := router.Group("/menu")
    menu.Use(handler.SessionCheck())
    {
        //物品一覧画面
        menu.GET("/gallery", handler.DisplayGalleryPage(Equips_info))   
        //所有物画面
        menu.GET("/mylist", handler.DisplayMylistPage(Equips_info))
        //パスワード変更画面
        menu.GET("/changepass", handler.DisplayChangepassPage)
        //パスワード変更
        menu.POST("/changepass", handler.ChangePassword)
        //物品借用申請
        menu.POST("/rent_apply", handler.RentApply)
        //物品返却申請
        menu.POST("/return_apply", handler.ReturnApply)
	}

    //Admin用のページ
    sec := router.Group("/sec")
    sec.Use(handler.SecCheck())
    {
        //管理画面
		sec.GET("/uadmin", handler.DisplayAdminPage(Equips_info))
        //貸出申請処理
        sec.POST("/rent_accept", handler.RentAccept)
        //ユーザリスト画面
        sec.GET("/list_user", handler.DisplayUserListPage)
        //ユーザ編集
        sec.POST("/edit_user", handler.EditUser)
        //カテゴリーリスト画面
        sec.GET("/list_classification", handler.DisplayCategoryListPage)	
        //カテゴリーリスト編集
        sec.POST("/edit_classification", handler.EditClassification)
        //ユーザ追加画面
        sec.GET("/add_user", handler.DisplayAddUserPage)
        //ユーザ追加
        sec.POST("/add_user", handler.AddUser)
        // 物品追加画面
        sec.GET("/add_equipment", handler.DisplayAddEquipmentPage)
        // 物品追加
        sec.POST("/add_equipment", handler.AddEquipment)
        // ジャンル追加画面
        sec.GET("/add_classification", handler.DisplayAddClassification)
        // ジャンル追加
        sec.POST("/add_classification", handler.AddClassicaioin)
        // 物品編集画面
        sec.GET("/list_equipment",handler.DisplayEquipmentListPage)
        // 物品編集
        sec.POST("/edit_equipment", handler.EditEquipment)
}

}


