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

// BatchProcessRecord 批量处理记录
type BatchProcessRecord struct {
	gorm.Model
	WordListIDs  string `json:"wordListIds"`
	MyVocabulary int    `json:"myVocabulary"`
	Vocabulary   int    `json:"vocabulary"`
	Difference   int    `json:"difference"`
}
