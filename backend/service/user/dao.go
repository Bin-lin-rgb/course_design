package user

import (
	"backend/common"
)

type UserInfo common.UserInfo

func (u *UserInfo) Create() error {
	result := db.Model(&UserInfo{}).Create(u)
	return result.Error
}

func (u *UserInfo) GetUserInfoByID() error {
	result := db.Model(&UserInfo{}).First(u)
	return result.Error
}

func (u *UserInfo) IsExist() error {
	result := db.Model(&UserInfo{}).Where("account = ? and password = ?", u.Account, u.Password).First(u)
	return result.Error
}

func (u *UserInfo) UpdatePassword() error {
	result := db.Model(&UserInfo{}).Where("account = ?", u.Account).Update("password", u.Password)
	return result.Error

	/*	if err := db.Where("account = ?", u.Account).First(u).Error; err != nil {
			// 如果找不到匹配记录，则返回错误

			return err
		}

		// 如果找到匹配记录，则更新记录
		result := db.Model(&UserInfo{}).Where("account = ?", u.Account).Update("password", u.Password)
		return result.Error*/
}
