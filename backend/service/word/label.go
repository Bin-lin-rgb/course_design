package word

import (
	"backend/utils"
)

var (
	response *utils.Response
)

type reqInsertRecord struct {
	Grade           string   `form:"grade" binding:"required"`           //年级
	Volume          string   `form:"volume" binding:"required"`          //册数
	Unit            string   `form:"unit" binding:"required"`            //单元
	Content         string   `form:"content" binding:"required"`         //内容
	Corpus          string   `form:"corpus" binding:"required"`          //教材语料
	ExpressProperty string   `form:"expressProperty" binding:"required"` //表现属性
	CulturalElem    []string `form:"culturalElem" binding:"required"`    //文化元素
	CarrierForm     []string `form:"carrierForm" binding:"required"`     //载体形式
}

type reqRecordList struct {
	Page     int64  `form:"page" binding:"required"`
	PageSize int64  `form:"pageSize" binding:"required"`
	Grade    string `form:"grade"`    //年级
	Volume   string `form:"volume"`   //册数
	Unit     string `form:"unit"`     //单元
	Content  string `form:"content"`  //内容
	Username string `form:"username"` //创建人
	IsLike   int64  `form:"isLike"`   //是否对创建人字段开启模糊查询，0->关闭、1->开启，默认关闭。
	IsDesc   int64  `form:"isDesc"`   //是否按照时间倒序排序
}

/*
func GetRecordList(c *gin.Context) {
	var (
		reqForm reqRecordList
		msg     string
		err     error
	)

	if err = c.ShouldBindQuery(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 获取可能的创建者 ID 数组
	creatorIds := make([]uint, 0)
	if reqForm.Username != "" {
		creator := UserInfo{}
		creators := make([]UserInfo, 0)
		creator.Username = reqForm.Username
		if reqForm.IsLike == 1 {
			creators, err = creator.GetUserInfoByUsername()
		} else {
			creators, err = creator.GetUserInfoByUsernameNoLike()
		}
		if err != nil {
			msg = "数据库查询失败"
			z.Error(fmt.Sprintf(msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		creatorIds = make([]uint, len(creators))
		for i, info := range creators {
			creatorIds[i] = info.ID
		}
	}

	textbook := Textbook{}
	page := int(reqForm.Page)
	pageSize := int(reqForm.PageSize)
	records, total, err := textbook.GetAllRecordByPage(page, pageSize, reqForm.Grade,
		reqForm.Volume, reqForm.Unit, reqForm.Content, reqForm.IsDesc, creatorIds)
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 获取 records 的 CreatorId 数组
	ids := make([]uint, len(records))
	for i, record := range records {
		ids[i] = record.CreatorId
	}
	// usernames 数组，与 CreatorId 数组 一一对应
	// usernames := make([]string, len(ids))
	// 通过 CreatorId 数组获取 []UserInfo
	// 不能直接在遍历 []UserInfo 得到 usernames 数组，因为 CreatorId 数组可能有重复的
	userinfo := UserInfo{}
	users, err := userinfo.GetUsernameByIds(ids)
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	usersMap := make(map[uint]string, len(users))
	for _, user := range users {
		usersMap[user.ID] = user.Username
	}

	type Item struct {
		ID              int            `json:"id"`
		Grade           string         `json:"grade"`           //年级
		Volume          string         `json:"volume"`          //册数
		Unit            string         `json:"unit"`            //单元
		Content         string         `json:"content"`         //内容
		Corpus          types.JSONText `json:"corpus"`          //教材语料
		ExpressProperty types.JSONText `json:"expressProperty"` //表现属性
		CulturalElem    types.JSONText `json:"culturalElem"`    //文化元素
		CarrierForm     types.JSONText `json:"carrierForm"`     //载体形式
		CreatedAt       string         `json:"createdAt"`       //创建时间
		CreatorId       uint           `json:"creatorId"`       //创建人ID
		Username        string         `json:"username"`        //创建者昵称
	}

	var items []Item
	for i, record := range records {
		corpus, _ := record.Corpus.MarshalJSON()
		expressProperty, _ := record.ExpressProperty.MarshalJSON()
		culturalElem, _ := record.CulturalElem.MarshalJSON()
		carrierForm, _ := record.CarrierForm.MarshalJSON()

		item := Item{
			ID:              int(record.Model.ID),
			Grade:           record.Grade,
			Volume:          record.Volume,
			Unit:            record.Unit,
			Content:         record.Content,
			Corpus:          corpus,
			ExpressProperty: expressProperty,
			CulturalElem:    culturalElem,
			CarrierForm:     carrierForm,
			CreatedAt:       record.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatorId:       record.CreatorId,
			Username:        usersMap[ids[i]],
		}
		items = append(items, item)
	}

	resp := struct {
		Items []Item `json:"items"`
		Total int64  `json:"total"`
	}{
		Items: items,
		Total: total,
	}

	response.Success(c, resp)

}

type reqRecord struct {
	Id int64 `form:"id" binding:"required"`
}

// GetRecord 获取一条记录
func GetRecord(c *gin.Context) {
	var (
		reqForm reqRecord
		msg     string
		err     error
	)

	if err = c.ShouldBindQuery(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	textbook := Textbook{}
	textbook.ID = uint(reqForm.Id)
	err = textbook.GetARecordById()
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	user := UserInfo{}
	user.ID = textbook.CreatorId
	err = user.GetUserInfoByID()
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	resp := struct {
		ID              uint           `json:"id"`
		CreatorId       uint           `json:"creatorId"`
		Username        string         `json:"username"`
		Grade           string         `json:"grade,omitempty"`           //年级
		Volume          string         `json:"volume,omitempty"`          //册数
		Unit            string         `json:"unit,omitempty"`            //单元
		Content         string         `json:"content,omitempty"`         //内容
		Corpus          types.JSONText `json:"corpus,omitempty"`          //教材语料
		ExpressProperty types.JSONText `json:"expressProperty,omitempty"` //表现属性
		CulturalElem    types.JSONText `json:"culturalElem,omitempty"`    //文化元素
		CarrierForm     types.JSONText `json:"carrierForm,omitempty"`     //载体形式
	}{
		ID:              textbook.ID,
		CreatorId:       textbook.ID,
		Username:        user.Username,
		Grade:           textbook.Grade,
		Volume:          textbook.Volume,
		Unit:            textbook.Unit,
		Content:         textbook.Content,
		Corpus:          textbook.Corpus,
		ExpressProperty: textbook.ExpressProperty,
		CulturalElem:    textbook.CulturalElem,
		CarrierForm:     textbook.CarrierForm,
	}

	response.Success(c, resp)

}

type reqModifyRecord struct {
	ID              uint     `json:"id" binding:"required"`
	Grade           string   `form:"grade" binding:"required"`           //年级
	Volume          string   `form:"volume" binding:"required"`          //册数
	Unit            string   `form:"unit" binding:"required"`            //单元
	Content         string   `form:"content" binding:"required"`         //内容
	Corpus          string   `form:"corpus" binding:"required"`          //教材语料
	ExpressProperty string   `form:"expressProperty" binding:"required"` //表现属性
	CulturalElem    []string `form:"culturalElem" binding:"required"`    //文化元素
	CarrierForm     []string `form:"carrierForm" binding:"required"`     //载体形式
}

// UpdateRecord 修改一条记录
func UpdateRecord(c *gin.Context) {
	var (
		reqForm reqModifyRecord
		msg     string
		err     error
	)

	if err = c.ShouldBind(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	corpus, _ := json.Marshal(reqForm.Corpus)
	expressProperty, _ := json.Marshal(reqForm.ExpressProperty)
	culturalElem, _ := json.Marshal(reqForm.CulturalElem)
	carrierForm, _ := json.Marshal(reqForm.CarrierForm)

	// 获取创建者Id
	temp, ok := c.Get("userID")
	if !ok {
		msg = "服务器出错，请稍后重试"
		z.Error("从 *gin.Context 获取用户 ID 失败")
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	// 只能是先转 int64 或 string
	userID, _ := temp.(int64)

	text := Textbook{}
	text.ID = reqForm.ID
	err = text.GetARecordById()
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	// 不是本人，无法修改
	if uint(userID) != text.CreatorId {
		msg = "抱歉，无权修改他人创建的语料"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	textbook := Textbook{
		Grade:           reqForm.Grade,
		Volume:          reqForm.Volume,
		Unit:            reqForm.Unit,
		Content:         reqForm.Content,
		Corpus:          corpus,
		ExpressProperty: expressProperty,
		CulturalElem:    culturalElem,
		CarrierForm:     carrierForm,
	}
	textbook.ID = reqForm.ID
	err = textbook.Update()
	if err != nil {
		msg = "数据库更新失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	response.Success(c, nil)
}

type reqDeleteRecord struct {
	Id int64 `form:"id" binding:"required"`
}

// DeleteRecord 删除一条记录
func DeleteRecord(c *gin.Context) {
	var (
		reqForm reqDeleteRecord
		msg     string
		err     error
	)

	if err = c.ShouldBindQuery(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 获取创建者Id
	temp, ok := c.Get("userID")
	if !ok {
		msg = "服务器出错，请稍后重试"
		z.Error("从 *gin.Context 获取用户 ID 失败")
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	// 只能是先转 int64 或 string
	userID, _ := temp.(int64)

	text := Textbook{}
	text.ID = uint(reqForm.Id)
	err = text.GetARecordById()
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	// 不是本人，无法删除
	if uint(userID) != text.CreatorId {
		msg = "抱歉，无权删除他人创建的语料"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	err = text.DeleteRecord()
	if err != nil {
		msg = "数据库删除失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	response.Success(c, nil)
	return
}
*/
