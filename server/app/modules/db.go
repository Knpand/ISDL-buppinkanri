package modules

import (
	//_　で使ってない場合のエラー回避
	//  "database/sql" 
	//  "log"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
    "app/constants"
	_ "net/http"
	_ "time"
    "strconv"
	_ "github.com/go-sql-driver/mysql"
)


/////////////////////////////////////////
//
// DB Connect用　db.Close()のベストな位置考えたい，main文に書きたくない
//
/////////////////////////////////////////
func GormConnect() *gorm.DB {
    DBMS := "mysql"
    USER := "light"
    PASS := "light"
    DBNAME := "IOU_db"
	PROTOCOL :="tcp(IOU_database:3306)"
    // MySQLだと文字コードの問題で"?parseTime=true"を末尾につける必要がある
    CONNECT := USER + ":" + PASS + "@"+PROTOCOL+"/" + DBNAME + "?charset=utf8&parseTime=true"
    db, err := gorm.Open(DBMS, CONNECT)
 
    if err != nil {
        panic(err.Error())
    }
    return db
}

/////////////////////////////////////////
//
// Get User Information in Admin
//
/////////////////////////////////////////

func Get_userinfo() []constants.Student {
    db := GormConnect()

    // defer db.Close()
    student := []constants.Student{} 
    // FindでDB名を指定して取得した後、orderで登録順に並び替え
    // db.Order("created_at desc").Find(&student)
    db.Find(&student)
    // fmt.Print("\n")
    // fmt.Print(student)
    // fmt.Print("\n")
    return student
}

/////////////////////////////////////////
//
// Update User Information in Admin
//
/////////////////////////////////////////

func Update_userinfo_admin(student constants.Student){
    db := GormConnect()
    defer db.Close()
    db.Table("students").Where("id = ?", student.Id).Update(&student)
}

/////////////////////////////////////////
//
// Delete User Information
//
/////////////////////////////////////////

func Delete_userinfo(student constants.Student){
    db := GormConnect()
    defer db.Close()
    db.Table("students").Where("id = ?", student.Id).Delete(&student)
}

/////////////////////////////////////////
//
// Get Class Information
//
/////////////////////////////////////////

func Get_classinfo() []constants.Classification {
    db := GormConnect()
    // defer db.Close()
    classification := []constants.Classification{} 
    db.Find(&classification)
    // fmt.Print("\n")
    // fmt.Print(classification)
    // fmt.Print("\n")
    return classification
}

/////////////////////////////////////////
//
// Update Class Information in Admin
//
/////////////////////////////////////////

func Update_classinfo_admin(classification constants.Classification){
    db := GormConnect()
    defer db.Close()
    db.Table("classifications").Where("id = ?", classification.Id).Update(&classification)
}

/////////////////////////////////////////
//
// Delete Class Information
//
/////////////////////////////////////////

func Delete_classinfo(classification constants.Classification){
    db := GormConnect()
    defer db.Close()
    db.Table("classifications").Where("id = ?", classification.Id).Delete(&classification)
}
/////////////////////////////////////////
//
// Update Equipments Information in Admin
//
/////////////////////////////////////////

func Update_equipinfo_admin(equipment constants.Equipment){
    db := GormConnect()
    defer db.Close()
    db.Table("equipments").Where("id = ?", equipment.Id).Update(&equipment)
}

/////////////////////////////////////////
//
// Insert Equipments Information
//
/////////////////////////////////////////

func Insert_equipinfo(name string, classification_id int){
    db := GormConnect()
    defer db.Close()
    new_data:=constants.Equipment{
        Name:name,
        Classifications_id:classification_id,
        State:"貸出可能",
    }
    db.Table("equipments").Create(&new_data)
}

/////////////////////////////////////////
//
// Insert Class Information
//
/////////////////////////////////////////

func Insert_classificationinfo(name string){
    db := GormConnect()
    defer db.Close()
    new_data:=constants.Classification{
        Name:name,
    }
    db.Table("classifications").Create(&new_data)
}

/////////////////////////////////////////
//
// Get User Password
//
/////////////////////////////////////////
func GetDBUserPassword(id string) string {
    db := GormConnect()
    var student constants.Student 
    db.Where("student_id = ?",id).First(&student)
    db.Close()
	fmt.Print(student)
    return student.Password
}

/////////////////////////////////////////
//
// Register User in GUI
//
/////////////////////////////////////////

func Register_user(c *gin.Context)[]error {
	// バリデーション処理
	student_id := c.PostForm("Student_id")
	name := c.PostForm("Name")
	password := c.PostForm("Password")
    email :=c.PostForm("Email")
    is_superuser:=c.PostForm("Is_superuser")
    superuser_bool,_:=strconv.ParseBool(is_superuser)
    fmt.Println(superuser_bool)
	return createUser(student_id,name, password,email,superuser_bool)
}

/////////////////////////////////////////
//
// Register User in DB
//
/////////////////////////////////////////

func createUser(student_id string,name string, password string,email string ,is_superuser bool) []error {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    db := GormConnect()
   
    // Insert処理
	student:= constants.Student{Student_id:student_id, Name:name,Email:email,Is_superuser:is_superuser,Password:string(hash)}
	// fmt.Print(student)
	//INSERT INTO students VALUES(Id,Student_i,Name,Password)
	//Createの変数にsをつけたテーブルが対象になる
	//Structの中の変数は最初の文字が大文字じゃないと参照されない
    if err := db.Create(&student).GetErrors(); err != nil {
		// fmt.Print(err)
        return err
    }
    defer db.Close()
    return nil
}

/////////////////////////////////////////
//
// Update User Information in DB
//
/////////////////////////////////////////

func Update_student(prop string , student_id string , newdata interface{}) {
    if prop == "password"{
        newdata, _ = bcrypt.GenerateFromPassword([]byte(newdata.(string)), bcrypt.DefaultCost)
    }
    db := GormConnect()
    defer db.Close()
    db.Table("students").Where("student_id = ?", student_id).Update(prop,newdata)
}

/////////////////////////////////////////
//
// Update equipments infomation
//
/////////////////////////////////////////

func Update_equipinfo(equip_info constants.Equipment){
    db := GormConnect()
    defer db.Close()
    db.Table("equipments").Where("id = ?", equip_info.Id).Update(&equip_info)
}   

/////////////////////////////////////////
//
// Delete Equipments infomation
//
/////////////////////////////////////////

func Delete_equipinfo(equip_info constants.Equipment){
    db := GormConnect()
    defer db.Close()
    db.Table("equipments").Where("id = ?", equip_info.Id).Delete(&equip_info)
}   

/////////////////////////////////////////
//
// Get Equipments Information
//
/////////////////////////////////////////

func Get_equipinfo() []constants.Equipment {
	db:=GormConnect()
	defer db.Close()
    var equipment []constants.Equipment
	// result:=db.Find(&equipment)
	db.Table("equipments").Find(&equipment)
	return equipment
}

/////////////////////////////////////////
//
// Get Equipments Class by Classid
//
/////////////////////////////////////////

func Getclassname(classid int) constants.Classification {
	db:=GormConnect()
	defer db.Close()
    var classification constants.Classification
	// result:=db.Find(&equipment)
	db.First(&classification, classid)
    // fmt.Println(classification)
	return classification
}

/////////////////////////////////////////
//
// Get Studentname by Studentid
//
/////////////////////////////////////////

func Select_Student_by_student_id(student_id string) constants.Student{
	db:=GormConnect()
	defer db.Close()
	var student constants.Student
	db.Table("students").Where("student_id = ?", student_id).First(&student)
	return student
}
