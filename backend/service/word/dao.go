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

type Wordbook2 struct {
	Wordbook
	Known string `json:"known"`
}

func GetWordListByWordArray(list []string) (wordList2 []Wordbook2, err error) {
	var wordList []Wordbook
	result := db.Model(&Wordbook{}).Select("id,word").Where("word in ?", list).Find(&wordList)
	for _, v := range wordList {
		wordList2 = append(wordList2, Wordbook2{v, ""})
	}
	return wordList2, result.Error
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

func (u *UserInfo) GetUserInfoByID() error {
	result := db.Model(&UserInfo{}).First(u)
	return result.Error
}

func (u *UserInfo) UpdateVocabulary() error {
	result := db.Model(&UserInfo{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"basic_vocabulary": u.BasicVocabulary,
	})
	return result.Error
}

func GetFourGradeAndVocabulary() (users []UserInfo, err error) {
	err = db.Model(&UserInfo{}).Select("id,four_grade,basic_vocabulary").Find(&users).Error
	return users, err
}

func GetSixGradeAndVocabulary() (users []UserInfo, err error) {
	err = db.Model(&UserInfo{}).Select("id,six_grade,basic_vocabulary").Find(&users).Error
	return users, err
}
