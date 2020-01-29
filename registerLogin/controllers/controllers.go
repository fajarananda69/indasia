package controllers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"indasia/registerLogin/config"
	"indasia/registerLogin/mail"
	"indasia/registerLogin/models"
	"log"
	"net/http"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
)

type Velidate struct {
	Email string `json:"email" orm:"size(128)"`
}

var ORM orm.Ormer

func init() {
	config.ConnectToDb()
	ORM = config.GetOrmObject()
}

func (r *Velidate) Verify() bool {
	var val bool
	query := fmt.Sprintf("SELECT * FROM login.users WHERE email = '%s' ", r.Email)
	err := ORM.Raw(query).QueryRow(&r)
	if err == nil {
		return true
	} else {
		return false
	}
	return val
}

func SetRegister(c *gin.Context) {
	var newUser models.Users
	c.BindJSON(&newUser)

	req := Velidate{
		Email: newUser.Email,
	}

	res := req.Verify()
	if res != true {
		email := newUser.Email

		token := models.GetTOTPToken(email)
		data := map[string]interface{}{
			"email":    newUser.Email,
			"password": newUser.Password,
			"username": newUser.UserName,
			"phone":    newUser.Phone,
			"image":    newUser.Image,
			"token":    token,
		}
		data1, _ := json.Marshal(data)
		clients := config.RedisConn()
		key := hex.EncodeToString([]byte(email))
		err := clients.Set(key, data1, 600*time.Second).Err()
		// err := clients.Set(key, string(data1), 0).Err()
		if err != nil {
			panic(err)
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": http.StatusInternalServerError, "error": "Email is exists"})
		} else {

			mail.SendMail(email, "http://localhost:3000/indasia/validateToken/"+key)

			c.JSON(http.StatusOK,
				gin.H{"status": http.StatusOK, "response": "check verification in your email"})
		}
	} else {
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Email is exist"})
	}

}

func saveRegister(data map[string]interface{}) models.ReturnData {
	// key := c.Param("key")
	// clients := config.RedisConn()
	// val, err := clients.Get(key).Result()
	// if err != nil {
	// 	panic(err)
	// }
	// var data map[string]interface{}

	// if err := json.Unmarshal([]byte(val), &data); err != nil {
	// 	panic(err)
	// }
	var response models.ReturnData

	email := data["email"].(string)
	password := models.Encrypt(data["password"].(string))
	username := data["username"].(string)
	phone := data["phone"].(string)
	image := data["image"].(string)

	req := Velidate{
		Email: email,
	}
	res := req.Verify()

	if res != true {
		query := fmt.Sprintf("INSERT INTO login.users (email,password,user_name,phone,image) VALUES ('%s','%s','%s','%s','%s') ", email, password, username, phone, image)
		_, errs := ORM.Raw(query).Exec()

		if errs == nil {
			// c.JSON(http.StatusOK, gin.H{
			// 	"status":    http.StatusOK,
			// 	"email":     email,
			// 	"user_name": username,
			// 	"phone":     phone,
			// 	"image":     image})
			mapRes := map[string]interface{}{
				"email":     email,
				"user_name": username,
				"phone":     phone,
				"image":     image,
			}
			response = models.ReturnData{Status: http.StatusOK, Data: mapRes, Message: "Sucsess"}
		} else {
			log.Panic()
			response = models.ReturnData{Status: http.StatusInternalServerError, Data: "", Message: "Register failed"}
		}
	} else {
		response = models.ReturnData{Status: http.StatusInternalServerError, Data: "", Message: "Email is Exist"}
	}
	return response
}

func Login(c *gin.Context) {
	var user models.Users
	var login models.Login
	c.BindJSON(&login)
	email := login.Email
	password := models.Encrypt(login.Password)
	fmt.Println("DAD ", models.Decrypt(password), password)
	query := fmt.Sprintf("SELECT * FROM login.users WHERE email = '%s' AND password = '%s' ", email, password)
	err := ORM.Raw(query).QueryRow(&user)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "users": &user})
	} else {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "error": "Failed login"})
	}
}

func ForgotPass(c *gin.Context) {
	var forgot models.Login
	c.BindJSON(&forgot)
	req := Velidate{
		Email: forgot.Email,
	}
	res := req.Verify()

	// email := forgot.Email
	// query := fmt.Sprintf("SELECT email FROM login.users WHERE email = '%s' ", email)
	// _, err := ORM.Raw(query).Exec()
	if res == true {
		m := mail.SendMail(forgot.Email, "http://localhost:3000/indasia/newpass")
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
	var user models.Users
	var newpass models.Login

	c.BindJSON(&newpass)

	req := Velidate{
		Email: newpass.Email,
	}
	res := req.Verify()
	if res == true {
		email := newpass.Email
		password := models.Encrypt(newpass.Password)
		query := fmt.Sprintf("UPDATE login.users SET password = '%s' WHERE email = '%s' ", password, email)
		err := ORM.Raw(query).QueryRow(&user)
		if err != nil {
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

func ValidateToken(c *gin.Context) {
	key := c.Param("key")
	token := c.Query("token")
	fmt.Println("Token", token)

	clients := config.RedisConn()
	val, err := clients.Get(key).Result()
	if err != nil {
		panic(err)
	}
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
