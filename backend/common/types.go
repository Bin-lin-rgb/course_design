package common

import (
	"gorm.io/gorm"
)

// UserInfo 用户信息
type UserInfo struct {
	gorm.Model
	Account   string `json:"account" gorm:"unique"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	FourGrade string `json:"fourGrade"`
	SixGrade  string `json:"sixGrade"`
}

// Wordbook 语料信息
type Wordbook struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Word string `json:"word"`
}
