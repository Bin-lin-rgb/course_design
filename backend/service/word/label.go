package word

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

var (
	response *utils.Response
)

func GetWordList1(c *gin.Context) {
	var (
		msg string
		err error
	)

	// 设置抽样参数
	// 从多少个单词抽
	totalWords := 10000
	// 分几层
	numLayers := 10
	// 抽出来多少个
	sampleSize := 40
	// 定义层次权重，采用指数衰减的方式
	weights := []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1, 0.05, 0.03, 0.02, 0.01}

	// 计算每个层次有多少单词数量 = 划分层级
	layerSizes := make([]int, numLayers)
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}
	remainingWords := totalWords
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(remainingWords) * weights[i] / totalWeight)
		remainingWords -= layerSizes[i]
	}

	// 每个层级应该取多少个单词
	randCount := 0
	shouldRand := make([]int, numLayers)
	for i := 0; i < numLayers; i++ {
		// 计算百分比
		tmpPer := weights[i] / totalWeight
		tmpSize := tmpPer * float64(sampleSize)
		shouldRand[i] = int(tmpSize)
		randCount += shouldRand[i]
	}

	if randCount < sampleSize {
		shouldRand[0] += sampleSize - randCount
	} else if len(shouldRand) > sampleSize {
		shouldRand[0] -= randCount - sampleSize
	}

	//	生成对应数量的随机数
	rand.Seed(time.Now().UnixNano())

	randomNumbers2 := make([]int, 0)
	curCount := 0
	for i := 0; i < numLayers; i++ {
		curCount += layerSizes[i]
		randomNumbers := generateRandomNumbers(shouldRand[i], curCount-layerSizes[i]+1, curCount)
		randomNumbers2 = append(randomNumbers2, randomNumbers...)
	}

	// 对 randomNumbers2 进行排序
	sort.Ints(randomNumbers2)

	list, err := GetWordListById(randomNumbers2)
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	response.Success(c, list)
	return
}

type reqWordListItem struct {
	Id    int    `json:"id" binding:"required"`
	Word  string `json:"word" binding:"required"`
	Known int    `json:"known" binding:"required"` //
}

type reqWordList struct {
	WordList []reqWordListItem `json:"wordList" binding:"required"`
}

// JudgeUserWordLevel 根据是否认识单词判断词汇等级，以推出新的单词列表
func JudgeUserWordLevel(c *gin.Context) {
	var (
		reqForm reqWordList
		msg     string
		err     error
	)
	if err = c.ShouldBindJSON(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 设置抽样参数
	// 从多少个单词抽
	totalWords := 10000
	// 分几层
	numLayers := 10
	// 抽出来多少个
	sampleSize := 40
	// 定义层次权重，采用指数衰减的方式
	weights := []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1, 0.05, 0.03, 0.02, 0.01}
	// 定义层次的得分权重
	knownWeights := []float64{0.1, 0.15, 0.15, 0.25, 0.25, 0.3, 0.3, 0.3, 0.35, 0.35}
	totalKnownWeights := 0.0
	for i := 0; i < numLayers; i++ {
		totalKnownWeights += knownWeights[i]
	}

	// 计算每个层次有多少单词数量 = 划分层级
	layerSizes := make([]int, numLayers)
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}
	remainingWords := totalWords
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(remainingWords) * weights[i] / totalWeight)
		remainingWords -= layerSizes[i]
	}

	// 每个层级应该取多少个单词
	randCount := 0
	shouldRand := make([]int, numLayers)
	for i := 0; i < numLayers; i++ {
		// 计算百分比
		tmpPer := weights[i] / totalWeight
		tmpSize := tmpPer * float64(sampleSize)
		shouldRand[i] = int(tmpSize)
		randCount += shouldRand[i]
	}

	if randCount < sampleSize {
		shouldRand[0] += sampleSize - randCount
	} else if len(shouldRand) > sampleSize {
		shouldRand[0] -= randCount - sampleSize
	}

	//每个层级的认识率
	knownPers := make([]float64, numLayers)
	curId := 0
	for i := 0; i < numLayers; i++ {
		knownCount := 0
		for ; curId < shouldRand[i]; curId++ {
			if reqForm.WordList[curId].Known == 1 {
				knownCount++
			}
		}
		// 计算该层级的认识率
		knownPer := float64(knownCount) / float64(shouldRand[i])
		knownPers[i] = knownPer
	}

	// 根据认识率的得分权重计算每个层级的得分
	knownScores := make([]float64, numLayers)
	for i := 0; i < numLayers; i++ {
		knownScores[i] = knownPers[i] * knownWeights[i] * totalKnownWeights * 100
	}

	response.Success(c, nil)
	return

}

func generateRandomNumbers(n, min, max int) []int {
	if max-min+1 < n {
		fmt.Println("无法生成指定数量的不重复随机数")
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	numbers := make([]int, n)
	used := make(map[int]bool)

	for i := 0; i < n; i++ {
		randomNumber := rand.Intn(max-min+1) + min
		for used[randomNumber] {
			randomNumber = rand.Intn(max-min+1) + min
		}
		numbers[i] = randomNumber
		used[randomNumber] = true
	}

	return numbers
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
*/
