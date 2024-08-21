package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	//ID        int64  `json:"userid"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	MobileNo string `json:"mobile_no"`
	UserName string `json:"username"`
	Password string `json:"passwd"`
}

type TempUsers struct {
	ID          int64     `json:"userid"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	MobileNo    string    `json:"mobile_no"`
	UserName    string    `json:"username"`
	Password    string    `json:"passwd"`
	Isactivated int64     `json:"isactivated"`
	MailToken   int64     `json:"mailtoken"`
	Createdat   time.Time `json:"created_at"`
}

type UserLogin struct {
	// gorm.Model
	// ID       int64  `json:"userid"`
	UserName string `json:"username"`
	PassWord string `json:"passwd"`
	Token    string `json:"token"`
}

type JsonUser struct {
	ID       int64  `json:"userid"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	MobileNo string `json:"mobile_no"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

type UserEmail struct {
	Email string `json:"email"`
}

//	Token    string `json:"token";sql:"-"`

type Token struct {
	UserId uint32
	jwt.StandardClaims
}

type Products struct {
	ProdId     uint32
	Prod_name  string
	Prod_desc  string
	Prod_qty   int64
	Prod_unit  string
	Prod_cost  float32
	Prod_sell  float32
	Created_at time.Time
	Updated_at time.Time
}
