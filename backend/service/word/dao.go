package word

import "backend/common"

type Wordbook common.Wordbook
type UserInfo common.UserInfo

func GetWordListById(Ids []int) (wordList []Wordbook, err error) {
	result := db.Model(&Wordbook{}).Select("id,word").Where("id in ?", Ids).Find(&wordList)
	return wordList, result.Error
}

/*
func (t *Textbook) GetAllRecordByPage(page, pageSize int, grade, volume, unit, content string, IsDesc int64, creatorIds []uint) (records []Textbook, total int64, err error) {
	result := db.Model(&Textbook{})
	if grade != "" {
		result.Where("grade LIKE ?", "%"+grade+"%")
	}
	if volume != "" {
		result.Where("volume LIKE ?", "%"+volume+"%")
	}
	if unit != "" {
		result.Where("unit LIKE ?", "%"+unit+"%")
	}
	if content != "" {
		result.Where("content LIKE ?", "%"+content+"%")
	}
	if len(creatorIds) != 0 {
		result.Where("creator_id IN (?)", creatorIds)
	}
	result.Count(&total)
	if IsDesc == 1 {
		result.Order("created_at desc")
	}
	result.Limit(pageSize).Offset(pageSize * (page - 1)).Find(&records)
	return records, total, result.Error
}
*/

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
