package word

import "backend/common"

type Wordbook common.Wordbook
type UserInfo common.UserInfo

func GetWordListById(Ids []int) (wordList []Wordbook, err error) {
	result := db.Model(&Wordbook{}).Select("id,word").Where("id in ?", Ids).Find(&wordList)
	return wordList, result.Error
}

func (w *Wordbook) GetWordListByWord() (wordList Wordbook, err error) {
	result := db.Model(&Wordbook{}).Select("id,word").Where("word = ?", w.Word).Find(&wordList).Limit(1)
	return wordList, result.Error
}

func GetWordListByWordArray(list []string) (wordList []Wordbook, err error) {
	result := db.Model(&Wordbook{}).Select("id,word").Where("word in ?", list).Find(&wordList)
	return wordList, result.Error
}

func (u *UserInfo) GetUserInfoByUsername() (users []UserInfo, err error) {
	result := db.Model(&UserInfo{}).Select("id,username").Where("username LIKE ?", "%"+u.Username+"%").Find(&users)
	return users, result.Error
}

func (u *UserInfo) GetUserInfoByUsernameNoLike() (users []UserInfo, err error) {
	result := db.Model(&UserInfo{}).Select("id,username").Where("username = ?", u.Username).Find(&users)
	return users, result.Error
}

func (u *UserInfo) GetUsernameByIds(ids []uint) (users []UserInfo, err error) {
	err = db.Select("id, username").Where("id IN (?)", ids).Distinct().Find(&users).Error
	return users, err
}
