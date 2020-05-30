package controllers

import (
	"encoding/json"
	"fmt"
	"indasia/registerLogin/helper"
	"indasia/registerLogin/models"
	"indasia/registerLogin/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// type CheckEmail struct {
// 	Email string `json:"email" orm:"size(128)"`
// }

// var ORM orm.Ormer

// func init() {
// 	config.ConnectToDb()
// 	ORM = config.GetOrmObject()
// }

// func (r *CheckEmail) CheckEmail() bool {
// 	var val bool
// 	query := fmt.Sprintf("SELECT * FROM login.users WHERE email = '%s' ", r.Email)
// 	err := ORM.Raw(query).QueryRow(&r)
// 	if err == nil {
// 		return true
// 	} else {
// 		return false
// 	}
// 	return val
// }

func SetRegister(c *gin.Context) {
	var newUser models.Users
	// clients := config.RedisConn()
	c.BindJSON(&newUser)

	req := service.CheckUser{
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	res := req.CheckEmail()
	pass := req.CheckPass()
	if res != true && pass == true {
		email := newUser.Email

		token := helper.GetTOTPToken(email)
		data := map[string]interface{}{
			"email":    newUser.Email,
			"password": newUser.Password,
			"username": newUser.UserName,
			"phone":    newUser.Phone,
			"image":    newUser.Image,
			"token":    token,
		}
		dataString, _ := json.Marshal(data)
		key := helper.ToMd5(email)
		err := service.SetRedis(key, dataString, 600)
		// err := clients.Set(key, data1, 600*time.Second).Err()
		// err := clients.Set(key, string(data1), 0).Err()
		if err != nil {
			panic(err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Email is exists"})
		} else {

			helper.SendMail(email, "http://localhost:3000/indasia/validateToken/"+key)

			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "response": "please check verification in your email"})
		}
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Email is exist or Password invalid"})
	}

}

func ValidateToken(c *gin.Context) {
	key := c.Param("key")
	token := c.Query("token")

	// clients := config.RedisConn()
	// val, err := clients.Get(key).Result()
	// if err != nil {
	// 	panic(err)
	// }

	val := service.GetRedis(key)
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(val), &data); err != nil {
		panic(err)
	}
	if token == data["token"] {
		save := saveRegister(data)
		c.JSON(http.StatusOK, gin.H{"status": save.Status, "users": save.Data, "message": save.Message})
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "invalid token"})
	}

}

func saveRegister(data map[string]interface{}) models.ReturnData {

	var response models.ReturnData

	// email := data["email"].(string)
	// password := helper.Encrypt2(data["password"].(string))
	// username := data["username"].(string)
	// phone := data["phone"].(string)
	// image := data["image"].(string)

	req := service.CheckUser{
		Email: data["email"].(string),
	}
	res := req.CheckEmail()

	if res != true {
		errs := service.InsertDB(data)
		// query := fmt.Sprintf("INSERT INTO login.users (email,password,username,phone,image) VALUES ('%s','%s','%s','%s','%s') ", email, password, username, phone, image)
		// _, errs := ORM.Raw(query).Exec()

		if errs == nil {
			mapRes := map[string]interface{}{
				"email":    data["email"].(string),
				"username": data["username"].(string),
				"phone":    data["phone"].(string),
				"image":    data["image"].(string),
			}
			response = models.ReturnData{Status: http.StatusOK, Data: mapRes, Message: "Sucsess"}
		} else {
			fmt.Println(errs)
			log.Panic()
			response = models.ReturnData{Status: http.StatusInternalServerError, Data: "", Message: "Register failed"}
		}
	} else {
		response = models.ReturnData{Status: http.StatusInternalServerError, Data: "", Message: "Email is Exist"}
	}
	return response
}

func Login(c *gin.Context) {
	// var user models.Users
	var login models.Login
	c.BindJSON(&login)
	// email := login.Email
	// password := helper.Encrypt2(login.Password)

	req := service.CheckUser{
		Email:    login.Email,
		Password: helper.Encrypt2(login.Password),
	}
	res := req.CheckLogin()
	// status := service.CheckLogin(email, password)

	switch {
	case strings.Contains(res, "active"):
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": "Success Login"})
	case strings.Contains(res, "pendding"):
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Account pennding"})
	default:
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Account stop"})
	}
	// query := fmt.Sprintf("SELECT * FROM login.users WHERE email = '%s' AND password = '%s' ", email, password)
	// err := ORM.Raw(query).QueryRow(&user)
	// if err == nil {
	// 	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": "Login Success"})
	// } else {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusInternalServerError,
	// 		gin.H{"status": http.StatusInternalServerError, "error": "Failed login"})
	// }
}

func ForgotPass(c *gin.Context) {
	var forgot models.Login
	c.BindJSON(&forgot)
	req := service.CheckUser{
		Email: forgot.Email,
	}
	res := req.CheckEmail()

	// email := forgot.Email
	// query := fmt.Sprintf("SELECT email FROM login.users WHERE email = '%s' ", email)
	// _, err := ORM.Raw(query).Exec()
	if res == true {
		m := helper.SendMail(forgot.Email, "http://localhost:3000/indasia/newpass")
		if m == true {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "response": "success send mail"})
		} else {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Failed sent email"})
		}
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Email not found"})
	}

}

func SetNewPass(c *gin.Context) {
	// var user models.Users
	var newpass models.Login

	c.BindJSON(&newpass)

	req1 := service.CheckUser{
		Email:    newpass.Email,
		Password: newpass.Password,
	}
	res := req1.CheckEmail()
	pass := req1.CheckPass()
	if res == true && pass == true {
		email := newpass.Email
		password := helper.Encrypt2(newpass.Password)

		err := service.UpdatePass(email, password)
		// query := fmt.Sprintf("UPDATE login.users SET password = '%s' WHERE email = '%s' ", password, email)
		// err := ORM.Raw(query).QueryRow(&user)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": "Success resset password"})
		} else {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Failed reset password"})
		}
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Email not found"})
	}

}
