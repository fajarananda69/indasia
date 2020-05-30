package models

// Users - Model for the uses table
type Users struct {
	UserId      int    `json:"user_id" orm:"auto"`
	Email       string `json:"email" orm:"size(50)"`
	Password    string `json:"password" orm:"size(100)"`
	UserName    string `json:"username" orm:"size(50)"`
	Phone       string `json:"phone" orm:"size(15)"`
	Image       string `json:"image"`
	Role        string `json:"role"`
	CodeReferal string `json:"code_referal"`
	CreatedAt   int    `json:"created_at"`
	LastLogin   int    `json:"last_login"`
	Status      string `json:"status" orm:"size(8)"`
}

type Login struct {
	Email    string `json:"email" orm:"size(50)"`
	Password string `json:"password" orm:"size(100)"`
}

type ReturnData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}
