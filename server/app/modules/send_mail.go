package modules

import (
	"log"
	"net/smtp"
    "github.com/joho/godotenv"
	"os"
)


func Send_mail(student_name string,equipment_name string){
	godotenv.Load(".env")
	subject:="貸出申請の通知"
	message:=student_name+"さんから"+equipment_name+"の貸出申請が届きました。\n http://172.20.11.231:80"
	// Set up authentication information.
	auth := smtp.PlainAuth(
			"",
			os.Getenv("IOU_MailAccount"), // foo@gmail.com
			os.Getenv("GMAIL_API_KEY"),
			"smtp.gmail.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			os.Getenv("IOU_MailAccount"), 
			[]string{os.Getenv("IOU_MailAccount")},
			[]byte("To: <recipient>@gmail.com\r\n" +
			"Subject:" + subject + "\r\n" +
			"\r\n" +message),
	)
	if err != nil {
			log.Fatal(err)
	}
}