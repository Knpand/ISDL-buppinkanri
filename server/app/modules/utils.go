package modules

import (
	//_　で使ってない場合のエラー回避
	//  "database/sql" 
	//  "log"
	// "github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"fmt"
	"app/constants"
	"strconv"
	"strings"
	_ "net/http"
	_  "time"
	_  "github.com/go-sql-driver/mysql" 
)

/////////////////////////////////////////
//
// Get Equip_info from Db（need remake)
//
/////////////////////////////////////////

func Update_EquipInfo() []constants.Equip_info {
	cnt:=0
	equipments:=Get_equipinfo()
	equips:= make([]constants.Equip_info, len(equipments))
	
	for _, equip := range equipments {
		classname:=Getclassname(equip.Classifications_id).Name
		equip_info:=constants.Equip_info{equip.Id,equip.Name,equip.User,equip.Classifications_id,classname,equip.State,equip.Remarks}
		equips[cnt]=equip_info
		cnt=cnt+1
	}
	return equips
}

/////////////////////////////////////////
//
// Get Studentinfo by Session（need remake)
//
/////////////////////////////////////////

func Get_studentinfo_bySession(ctx *gin.Context ,prop string) interface{}{
	var setprop interface{}
	setprop, _ = ctx.Get(prop)
	return setprop
}

/////////////////////////////////////////
//
// Get Studentinfo by Student_id for Session
//
/////////////////////////////////////////


func Fetch_studentinfo_byID(student_id string) constants.Student{
	student_info:=Select_Student_by_student_id(student_id)
	return student_info
}


////////////////////////////////////////
//
// Check Form(Add User)
//
/////////////////////////////////////////

func Isexists_userform(ctx *gin.Context) (constants.Isexists_userform , bool){
	isexists_userform :=constants.Isexists_userform{true,true,true,true}
	isfull_form:=true
	a:=ctx.PostForm("Name")
	if a == ""{
		isexists_userform.Name=false
		isfull_form=false
	}
	a=ctx.PostForm("Student_id")
	if a== ""{
		isexists_userform.Student_id=false
		isfull_form=false
	}
	a=ctx.PostForm("Email")
	if a =="" {
		isexists_userform.Email=false
		isfull_form=false
	}
	a=ctx.PostForm("Password")
	if a==""{
		isexists_userform.Password=false
		isfull_form=false
	}
	return isexists_userform,isfull_form
}

////////////////////////////////////////
//
// Check Form(Add Classifications)
//
/////////////////////////////////////////

func Isexists_classificationsform(ctx *gin.Context) (bool , bool){
	classifications_form:=true
	isfull_form:=true
	a:=ctx.PostForm("Name")
	if a == ""{
		classifications_form=false
		isfull_form=false
	}
	return classifications_form,isfull_form
}

////////////////////////////////////////
//
// Check Form(Add Equipment)
//
/////////////////////////////////////////

func Isexists_equipform(ctx *gin.Context) (constants.Isexists_equipform , bool){
	isexists_equipform :=constants.Isexists_equipform{true,true}
	isfull_form:=true
	a:=ctx.PostForm("Name")
	if a == ""{
		isexists_equipform.Name=false
		isfull_form=false
	}
	b,_:=strconv.Atoi(ctx.PostForm("Classifications_id"))
	if b==0{
		isexists_equipform.Classifications_id=false
		isfull_form=false
	}
	
	return isexists_equipform,isfull_form
}

/////////////////////////////////////////
//
// Apply Equipments 
//
/////////////////////////////////////////

func Rent_apply(rent_info constants.Equipment){
	Update_equipinfo(rent_info)
	//ここにページリロード処理書きたい
}

func Return_apply(rent_info constants.Equipment){
	Update_equipinfo(rent_info)
	//ここにページリロード処理書きたい
}

/////////////////////////////////////////
//
// Accept Rent Request
//
/////////////////////////////////////////
func Accept_request(rent_info constants.Equipment){
	Update_equipinfo(rent_info)
}


/////////////////////////////////////////
//
// Change Password
//
/////////////////////////////////////////

func Change_passwords(student_id string, original_pass string, new_pass1 string, new_pass2 string)(bool,bool){
	db_pass:=GetDBUserPassword(student_id)
	var isEqual_originpass bool
	var isEqual_newpass bool
	if err := CompareHashAndPassword(db_pass, original_pass); err != nil {
		fmt.Println("not eq pass")
		isEqual_originpass=false
	} else {
		isEqual_originpass=true
		}
	if new_pass1 == new_pass2 && new_pass1 !=""{
		isEqual_newpass=true
	}else{
		isEqual_newpass=false
	}
	if isEqual_originpass && isEqual_newpass {
		Update_student("password",student_id,new_pass1)
	}
	return isEqual_originpass,isEqual_newpass
}

/////////////////////////////////////////
//
// Compare Hash with Password (パッケージ分けたい)
//
/////////////////////////////////////////
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

/////////////////////////////////////////
//
// Create Property about Showing Page
//
/////////////////////////////////////////
func Create_Pageinfo(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


func Hash_password(password string) string {
	new_password, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(new_password)
}

func Get_RentalApplyEquipment(ctx *gin.Context) constants.Equipment {
	session:=sessions.Default(ctx)
	equipment_id,_:=strconv.Atoi(ctx.PostForm("Equipment_id"))
	equipment_name:=ctx.PostForm("Name")
	classifications_id,_:=strconv.Atoi(ctx.PostForm("Classifications_id"))
	remarks:=ctx.PostForm("Remarks")
	studentname:=session.Get("Name").(string)

	rental_equipment := constants.Equipment{
		Id:equipment_id,
		Name:equipment_name,
		User:studentname,
		Classifications_id:classifications_id,
		State:"申請中",
		Remarks:remarks}

	return rental_equipment
}

func Get_ReturnApplyEquipment(ctx *gin.Context) constants.Equipment {
	
	equipment_id,_:=strconv.Atoi(ctx.PostForm("Equipment_id"))
	equipment_name:=ctx.PostForm("Name")
	classifications_id,_:=strconv.Atoi(ctx.PostForm("Classifications_id"))
	remarks:=" "
	studentname:=" "

	rental_equipment := constants.Equipment{
		Id:equipment_id,
		Name:equipment_name,
		User:studentname,
		Classifications_id:classifications_id,
		State:"貸出可能",
		Remarks:remarks}

	return rental_equipment
}

func Get_RentalAcceptEquipment(ctx *gin.Context) constants.Equipment {
	var state string
	var remarks string
	strIsAccept:=ctx.PostForm("IsAccept")
	isAccept,_:=strconv.ParseBool(strIsAccept)
	equipment_name:=ctx.PostForm("Name")
	equipment_id,_:=strconv.Atoi(ctx.PostForm("Equipment_id"))
	classifications_id,_:=strconv.Atoi(ctx.PostForm("Classifications_id")) 
	studentname:=ctx.PostForm("User")
	 if isAccept{
		state="貸出中"
		remarks=ctx.PostForm("Remarks")
	}else{
		state="貸出可能"
		remarks=" "
		studentname=" "
	}
   
	equipment := constants.Equipment{Id:equipment_id,
		Name:equipment_name,
		User:studentname,
		Classifications_id:classifications_id,
		State:state,
		Remarks:remarks}
	return equipment
}

func Get_EditUser(ctx *gin.Context) constants.Student {
	id,_ := strconv.Atoi(ctx.PostForm("Id"))
	student_id := ctx.PostForm("Student_id")
	email := ctx.PostForm("Email")
	name := ctx.PostForm("Name")
	tmp_is_superuser := ctx.PostForm("Is_superuser")
	var is_superuser bool
	if tmp_is_superuser == "false" {
		is_superuser = false
	} else {
		is_superuser = true
	}
	password := ctx.PostForm("Password")
	if password == "" {
		password=ctx.PostForm("Pre_Password")
	} else { 
		password=Hash_password(password)
		fmt.Print(password)
	}
	student := constants.Student{Id:id,
		Password:password,
		Student_id:student_id,
		Name:name,
		Email:email,
		Is_superuser:is_superuser}

	return student
}

func Get_EditEquipment(ctx *gin.Context) constants.Equipment {
	id,_ := strconv.Atoi(ctx.PostForm("Id"))
	name := ctx.PostForm("Name")
	user := ctx.PostForm("User")


	t_classifications_id := ctx.PostForm("Classifications_id")
	classifications_id,_ := strconv.Atoi(t_classifications_id[:strings.Index(t_classifications_id, "：")])
	state := ctx.PostForm("State")
	remarks := ctx.PostForm("Remarks")
	equipment := constants.Equipment{Id:id,
									Name:name,
									User:user,
									Classifications_id:classifications_id,
									State:state,
									Remarks:remarks}
	return equipment
}