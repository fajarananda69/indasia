package service

import (
	"fmt"
	"indasia/registerLogin/config"
	"indasia/registerLogin/helper"
	"indasia/registerLogin/models"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type CheckUser struct {
	Email    string
	Password string
	Token    string
	Status   string
}

var ORM orm.Ormer

func init() {
	config.ConnectToDb()
	ORM = config.GetOrmObject()
}

func (r *CheckUser) CheckEmail() bool {
	var val bool
	var user models.Login
	query := fmt.Sprintf("SELECT email FROM login.users WHERE email = '%s' ", r.Email)
	err := ORM.Raw(query).QueryRow(&user)
	if err == nil {
		return true
	} else {
		return false
	}
	return val
}

func (r *CheckUser) CheckPass() bool {
	var val bool
	regEmail, _ := regexp.Compile(`(^[A-z0-9.]+)`)
	matchEmail := regEmail.FindString(r.Email)

	if len(r.Password) >= 8 || len(r.Password) <= 15 {
		// regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
		regex1 := regexp.MustCompile(`[a-z]+`).MatchString
		regex2 := regexp.MustCompile(`[A-Z]+`).MatchString
		regex3 := regexp.MustCompile(`[0-9]+`).MatchString
		regexE := regexp.MustCompile(matchEmail).MatchString

		if regex1(r.Password) && regex2(r.Password) && regex3(r.Password) && !regexE(r.Password) {
			val = true
		} else {
			val = false
		}
	} else {
		val = false
	}
	return val
}

func InsertDB(data map[string]interface{}) error {
	createdAt := time.Now().Unix()
	email := data["email"].(string)
	password := helper.Encrypt2(data["password"].(string))
	username := data["username"].(string)
	phone := data["phone"].(string)
	image := data["image"].(string)
	role := "user"
	status := "active"
	codeReferal := getCodeReferal(email)

	query := fmt.Sprintf("INSERT INTO login.users (email,password,username,phone,image,created_at,role,code_referal,status_user) VALUES ('%s','%s','%s','%s','%s','%d','%s','%s','%s') ", email, password, username, phone, image, createdAt, role, codeReferal, status)
	_, errs := ORM.Raw(query).Exec()

	return errs
}

func (r *CheckUser) CheckLogin() string {
	var users models.Users

	query := fmt.Sprintf("SELECT status FROM login.users WHERE email = '%s' AND password = '%s' ", r.Email, r.Password)
	err := ORM.Raw(query).QueryRow(&users)

	status := fmt.Sprint(users.Status)
	if err != nil {
		fmt.Println("Email Password not found")
	} else {
		switch {
		case strings.Contains(status, "active"):
			lastLogin(r.Email)
		case strings.Contains(status, "pendding"):
			fmt.Println("Account pendding")
		default:
			fmt.Println("Account stop")
		}

	}
	return users.Status
}

func lastLogin(email string) {
	// var users models.Users
	lastLogin := time.Now().Unix()
	query := fmt.Sprintf("UPDATE login.users SET last_login = '%d' WHERE email = '%s' ", lastLogin, email)
	_, err := ORM.Raw(query).Exec()
	if err != nil {
		fmt.Println("failed set last Login")
	}
}

// func UpdatePass(email string, password string) error {

// 	query := fmt.Sprintf("UPDATE login.users SET password = '%s' WHERE email = '%s' ", password, email)
// 	_, err := ORM.Raw(query).Exec()
// 	return err
// }

func getCodeReferal(email string) string {
	code := helper.ToMd5(email)
	referal := code[len(code)-6:]
	return referal
}

func SetRedis(key string, data interface{}, sec time.Duration) error {
	clients := config.RedisConn()
	err := clients.Set(key, data, sec*time.Second).Err()
	if err != nil {
		fmt.Println("Failed send to redis")
	}
	return err
}

func GetRedis(key string) string {
	clients := config.RedisConn()
	val, err := clients.Get(key).Result()
	if err != nil {
		panic(err)
	}
	return val
}
