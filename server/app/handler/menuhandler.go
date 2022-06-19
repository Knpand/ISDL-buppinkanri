package handler

import (
	"github.com/gin-gonic/gin"
	"app/modules"
	"net/http"
	"github.com/gin-contrib/sessions"
	"app/constants"
)

func DisplayGalleryPage(Equips_info []constants.Equip_info)func(*gin.Context) {
	Equips_info=modules.Update_EquipInfo()  //物品情報の更新
	return func(ctx *gin.Context){
		session:=sessions.Default(ctx)
		root_links:=[]constants.Root_link{}
		page_info :=constants.Page_info{"物品一覧","/menu/gallery",root_links,"Equipments",session.Get("Is_superuser").(bool)}
		ctx.HTML(http.StatusOK, "gallery.html", gin.H{"page_info":page_info,"details":Equips_info})
	}

}

func DisplayMylistPage(Equips_info []constants.Equip_info)func(*gin.Context){
	Equips_info=modules.Update_EquipInfo()  //物品情報の更新
	return func(ctx *gin.Context){
		session:=sessions.Default(ctx)
		root_links:=[]constants.Root_link{}
		page_info :=constants.Page_info{"Mylist","/menu/mylist",root_links,"My List",session.Get("Is_superuser").(bool)}
		username:=session.Get("Name").(string)
		ctx.HTML(200, "mylist.html", gin.H{"page_info":page_info,"username":username,"details":Equips_info})
	}
}

func DisplayChangepassPage(ctx *gin.Context){
	session:=sessions.Default(ctx)
	root_links:=[]constants.Root_link{}
	page_info :=constants.Page_info{"パスワード変更","/menu/changepass",root_links,"Change Password",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "changepass.html", gin.H{"page_info":page_info,"isGet":true,"isEqual_originpass":true,"isEqual_newpass":true})
}

func ChangePassword(ctx *gin.Context){
	session:=sessions.Default(ctx)
	current_student_id:=session.Get("Student_id").(string)
	original_pass:=ctx.PostForm("Original_pass")
	new_pass1:=ctx.PostForm("New_pass1")
	new_pass2:=ctx.PostForm("New_pass2")
	isEqual_originpass,isEqual_newpass:=modules.Change_passwords(current_student_id, original_pass,new_pass1,new_pass2)
	root_links:=[]constants.Root_link{}
	page_info :=constants.Page_info{"パスワード変更","/menu/changepass",root_links,"Change Password",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "changepass.html", gin.H{"page_info":page_info,"isGet":false,"isEqual_originpass":isEqual_originpass,"isEqual_newpass":isEqual_newpass})
}

func RentApply(ctx *gin.Context){
	rental_equipment := modules.Get_RentalApplyEquipment(ctx)
	modules.Rent_apply(rental_equipment)
	ctx.Redirect(302, "/menu/gallery")
}

func ReturnApply(ctx *gin.Context){
	return_equipment := modules.Get_ReturnApplyEquipment(ctx)
	modules.Return_apply(return_equipment)  
	ctx.Redirect(302, "/menu/mylist")
}