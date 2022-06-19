package handler

import (
	"github.com/gin-gonic/gin"
	"app/modules"
	"github.com/gin-contrib/sessions"
	"log"
	"strconv"
	"app/constants"
)

func DisplayAdminPage(Equips_info []constants.Equip_info)func(*gin.Context){
	Equips_info=modules.Update_EquipInfo()  //物品情報の更新
	return func (ctx *gin.Context) {
		session:=sessions.Default(ctx)
		root_links:=[]constants.Root_link{}
		page_info :=constants.Page_info{"Admin","/sec/uadmin",root_links,"Admin",session.Get("Is_superuser").(bool)}
		ctx.HTML(200, "admin.html", gin.H{"page_info":page_info,"details":Equips_info})
	}
}

func RentAccept(ctx *gin.Context) {
	rent_equipment:=modules.Get_RentalAcceptEquipment(ctx)
	modules.Accept_request(rent_equipment)
	ctx.Redirect(302, "/sec/uadmin")
}

func DisplayUserListPage(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"list_user","/sec/list_user",root_links,"User List",session.Get("Is_superuser").(bool)}
	students := modules.Get_userinfo()
	ctx.HTML(200, "list_user.html", gin.H{"page_info":page_info,  "students": students})
}

func DisplayCategoryListPage(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	classifications := modules.Get_classinfo()
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"list_classification","/sec/list_classification",root_links,"Class List",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "list_classification.html", gin.H{"page_info":page_info, "classifications": classifications})
}

func EditClassification(ctx *gin.Context) {
	id,_ := strconv.Atoi(ctx.PostForm("Id"))
	name := ctx.PostForm("Name")
	classification := constants.Classification{id,name}
	delete :=ctx.PostForm("Isdelete")
	isdelete,_:=strconv.ParseBool(delete)
	change :=ctx.PostForm("Ischange")
	ischange,_:=strconv.ParseBool(change)
	if isdelete {
		modules.Delete_classinfo(classification) 
	}
	if ischange {
		 modules.Update_classinfo_admin(classification)
	}
	ctx.Redirect(302, "/sec/list_classification")
}

func EditUser(ctx *gin.Context) {
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
}

func DisplayAddUserPage(ctx *gin.Context){
	session:=sessions.Default(ctx)
	isexists_userform:=constants.Isexists_userform{true,true,true,true}
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_user","/sec/add_user",root_links,"Add User",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "add_user.html", gin.H{"page_info":page_info, "isexists_userform":isexists_userform})
}

func AddUser(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_user","/sec/add_user",root_links,"Add User",session.Get("Is_superuser").(bool)}
	isexists_userform,isfull_form :=modules.Isexists_userform(ctx)
	if isfull_form{
	if err :=  modules.Register_user(ctx); err != nil {
		log.Print(err)
	}
} 
ctx.HTML(200, "add_user.html", gin.H{"page_info":page_info, "isfull_form":isfull_form,"isexists_userform": isexists_userform})
}

func DisplayAddEquipmentPage(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	classifications := modules.Get_classinfo()
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_equipment","/sec/add_equipment",root_links,"Add Equipment",session.Get("Is_superuser").(bool)}
	isexists_equipform:=constants.Isexists_equipform{true,true}

	ctx.HTML(200, "add_equipment.html", gin.H{"page_info":page_info, "classifications": classifications,"isexists_equipform": isexists_equipform})

}

func AddEquipment(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_equipment","/sec/add_equipment",root_links,"Add Equipment",session.Get("Is_superuser").(bool)}
	name := ctx.PostForm("Name")
	classification_id,_ := strconv.Atoi(ctx.PostForm("Classifications_id"))

	isexists_equipform,isfull_form :=modules.Isexists_equipform(ctx)
	if isfull_form{
	modules.Insert_equipinfo(name,classification_id)
	}
	classifications := modules.Get_classinfo()

	ctx.HTML(200, "add_equipment.html", gin.H{"page_info":page_info, "classifications": classifications,"isfull_form":isfull_form,"isexists_equipform": isexists_equipform})

}

func DisplayAddClassification(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_classification","/sec/add_classification",root_links,"Add Class",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "add_classification.html", gin.H{"page_info":page_info,"isexists_classificationsform":true})
}

func AddClassicaioin(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	name := ctx.PostForm("Name")
	isexists_classificationsform,isfull_form:=modules.Isexists_classificationsform(ctx)
	if isfull_form {
	   modules.Insert_classificationinfo(name) 
	}
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"add_classification","/sec/add_classification",root_links,"Add Class",session.Get("Is_superuser").(bool)}

	ctx.HTML(200, "add_classification.html", gin.H{"page_info":page_info,"isexists_classificationsform":isexists_classificationsform})
}

func DisplayEquipmentListPage(ctx *gin.Context) {
	session:=sessions.Default(ctx)
	Equips_info:=modules.Update_EquipInfo()  //物品情報の更新
	classifications := modules.Get_classinfo()
	root_links:=[]constants.Root_link{constants.Root_link{"/sec/uadmin","Admin"}}
	page_info :=constants.Page_info{"list_equipment","/sec/list_equipment",root_links,"Equipment List",session.Get("Is_superuser").(bool)}
	ctx.HTML(200, "list_equipment.html", gin.H{"page_info":page_info,"details":Equips_info, "classifications": classifications})
}

func EditEquipment(ctx *gin.Context) {
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
	ctx.Redirect(302, "/sec/list_equipment")
}