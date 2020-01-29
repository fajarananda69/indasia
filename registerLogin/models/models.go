package models

// Users - Model for the uses table
type Users struct {
	UserId      int    `json:"user_id" orm:"auto"`
	Email       string `json:"email" orm:"size(128)"`
	Password    string `json:"password" orm:"size(64)"`
	UserName    string `json:"user_name" orm:"size(32)"`
	Phone       string `json:"phone" orm:"size(12)"`
	Image       string `json:"image"`
	CodeReferal int    `json:"code_referal" `
}

type Login struct {
	Password string `json:"password" orm:"size(64)"`
	Email    string `json:"email" orm:"size(128)"`
}

type ReturnData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}
