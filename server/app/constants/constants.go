package constants


type Student struct {
	Id		int  `gorm:"primaryKey`
	Student_id  string `gorm:"type:varchar(100)`
	Name	string  `gorm:"type:varchar(100)`
    Email string    `gorm:"type:varchar(100)`
    Is_superuser bool
	Password	string   `gorm:"type:varchar(128)`
}

type Rootlinks struct {
    Link string
    Name string
  }

type Equip_info struct{
    Id int
    Name string
	User string
	Classifications_id int
    Category string
    State string
    Remarks string
}

type Equipment struct {
	Id int  `gorm:"primaryKey`
	Name	string  `gorm:"type:varchar(50)`
	User string
	Classifications_id int 
	State string
	Remarks string
}



type Classification struct{
	Id int  `gorm:"primaryKey`
	Name	string  `gorm:"type:varchar(50)`
}

type Rent_info struct{
	Student Student
	Equipment Equipment
}

type Page_info struct{
	Pagename  string
	Currentlink string
	Rootlinks  []Root_link
	Locate  string
	IsAdmin bool
}

type Root_link struct{
	Link string
	Name string
}

type Isexists_userform struct{
	Name bool
	Student_id bool
	Email bool
	Password bool
}

type Isexists_equipform struct{
	Name bool
	Classifications_id bool
}