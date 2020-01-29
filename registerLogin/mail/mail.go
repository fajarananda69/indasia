package mail

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendMail(email string, url string) bool {
	mail := false
	e := godotenv.Load("registerLogin/conf.env")
	if e != nil {
		fmt.Print("SSSS", e)
	}

	configHost := os.Getenv("CONFIG_SMTP_HOST")
	configPort, _ := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	configEmail := os.Getenv("CONFIG_EMAIL")
	configPassword := os.Getenv("CONFIG_PASSWORD")

	// key := hex.EncodeToString([]byte(email))
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", configEmail)
	mailer.SetHeader("To", email)
	// mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", "Hello, <b>have a nice day</b> <a href='"+url+"'>link login klik here</a>")
	// mailer.Attach("./sample.png")

	dialer := gomail.NewPlainDialer(
		configHost,
		configPort,
		configEmail,
		configPassword,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
		fmt.Println("ERR ", err)
		mail = false
	} else {
		log.Println("Mail sent!")
		mail = true
	}

	return mail
}
