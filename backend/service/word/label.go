package word

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

var UploadFile = "./upload"

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
	totalWords := 8000
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
	//remainingWords := totalWords
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(totalWords) * weights[i] / totalWeight)
		//remainingWords -= layerSizes[i]
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
	if len(reqForm.WordList) < 3 {
		msg = "参数太少，请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数太少啦:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 设置抽样参数
	numLayers := 6
	totalWords := reqForm.WordList[len(reqForm.WordList)-1].Id

	// 定义层次权重，采用指数衰减的方式
	weights := []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1}
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}

	// 定义层次的得分权重
	knownWeights := []float64{1.0, 1.3, 1.6, 1.9, 2.2, 2.5}
	totalKnownWeights := 0.0
	for i := 0; i < numLayers; i++ {
		totalKnownWeights += knownWeights[i]
	}

	// 定义层次系数
	//levelCoefficient := []float64{0.95, 0.8, 0.85, 0.8, 0.75, 0.7}

	// 计算每个层次的分界线 = 划分层级
	boundaries := make([]int, numLayers+1)
	boundaries[0] = 1
	for i := 1; i <= numLayers; i++ {
		boundary := int(math.Round(float64(totalWords) * (weights[i-1] / totalWeight)))
		boundaries[i] = boundaries[i-1] + boundary
	}

	knownPers := make([]float64, numLayers)
	totalPer := 0.0
	curId := 0
	length := len(reqForm.WordList)

	for i := 0; i < numLayers; i++ {
		tmp := curId
		knownCount := 0
		layerBoundaryStart := boundaries[i]
		layerBoundaryEnd := boundaries[i+1]

		for ; curId < length && reqForm.WordList[curId].Id < layerBoundaryEnd; curId++ {
			if reqForm.WordList[curId].Id > layerBoundaryStart && reqForm.WordList[curId].Known == 1 {
				knownCount++
			}
		}

		// 计算该层级的认识率
		if curId > 0 {
			knownPer := float64(knownCount) / float64(curId-tmp)
			knownPers[i] = knownPer
			totalPer += knownPer
		} else {
			knownPers[i] = 0.0
		}
	}

	totalPer = totalPer / float64(numLayers)
	// 根据层级准确率计算下次抽样的数量
	sampleSize := 0
	// 从多少个单词抽
	wordsRange := 5000
	// weights: 定义层次权重，采用指数衰减的方式

	if totalPer < 0 || totalPer > 1 {
		msg = "计算的认识率不合法"
		z.Error(fmt.Sprintf("%s:%s,计算的认识率不合法:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	if totalPer < 0.3 {
		sampleSize = 30
		wordsRange = 5000
		weights = []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1}
	} else if totalPer < 0.5 {
		sampleSize = 60
		wordsRange = 15000
		weights = []float64{1.0, 0.8, 0.7, 0.6, 0.5, 0.4}
	} else if totalPer < 0.8 {
		sampleSize = 80
		wordsRange = 25000
		weights = []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5}
	} else {
		sampleSize = 100
		wordsRange = 41000
		weights = []float64{1.0, 0.9, 0.9, 0.9, 0.8, 0.7}
	}

	//根据层级准确率计算下次抽样的单词列表
	// 计算每个层次有多少单词数量 = 划分层级
	layerSizes := make([]int, numLayers)
	totalWeight2 := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight2 += weights[i]
	}
	//remainingWords := wordsRange
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(wordsRange) * weights[i] / totalWeight2)
		//remainingWords -= layerSizes[i]
	}

	// 每个层级应该取多少个单词
	randCount := 0
	shouldRand := make([]int, numLayers)
	for i := 0; i < numLayers; i++ {
		// 计算百分比
		tmpPer := weights[i] / totalWeight2
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

func CalculateVocabulary(c *gin.Context) {
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
	if len(reqForm.WordList) < 3 {
		msg = "参数太少，请求不合法"
		z.Error(fmt.Sprintf("%s:%s,请求的参数太少啦:%+v", msg, err, reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	//var list []Wordbook2
	list := make([]Wordbook2, len(reqForm.WordList))
	// 将 reqForm.WordList 转换成 workbook2
	for i := 0; i < len(reqForm.WordList); i++ {
		list[i].ID = uint(reqForm.WordList[i].Id)
		list[i].Word = reqForm.WordList[i].Word
		list[i].Known = strconv.Itoa(reqForm.WordList[i].Known)
	}

	sortByID(list)

	vo := EstimateVocabulary(list)

	userID, _ := c.Get("userID")
	user := UserInfo{
		Model: gorm.Model{ID: uint(userID.(int64))},
	}
	err = user.GetUserInfoByID()
	if err != nil {
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	var array []int
	err = json.Unmarshal([]byte(user.BasicVocabulary), &array)

	array = append(array, vo)

	bytes, err := json.Marshal(array)
	if err != nil {
		msg = "内部服务器错误"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	user.BasicVocabulary = string(bytes)
	err = user.UpdateVocabulary()
	if err != nil {
		msg = "数据库更新失败"
		z.Error(fmt.Sprintf(msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	response.Success(c, vo)
	return
}

// 定义排序函数
func sortByID(wordbooks []Wordbook2) {
	// 自定义比较函数
	less := func(i, j int) bool {
		return wordbooks[i].ID < wordbooks[j].ID
	}

	// 使用排序函数对数组进行排序
	sort.Slice(wordbooks, less)
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

type resBatchProcessItem struct {
	Round            int    `json:"round"`
	TestVocabulary   int    `json:"test_vocabulary"`
	PreplyVocabulary string `json:"preply_vocabulary"`
	CalculatedWords  int    `json:"calculated_words"`
}

func BatchProcess(c *gin.Context) {
	var (
		msg string
		err error
	)

	file, err := c.FormFile("file")
	if err != nil {
		msg = "文件上传失败"
		z.Error(fmt.Sprintf("%s:%s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	//获取文件名
	fileName := file.Filename

	if file.Size > 10*1024*1024 {
		msg = "文件过大"
		z.Error(fmt.Sprintf("文件过大:%s", msg))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	newUUID, _ := uuid.NewUUID()

	// 文件后缀
	fileSuffix := path.Ext(fileName)
	// 路径
	dst := path.Join(UploadFile, newUUID.String()+fileSuffix)

	// 先判断目录是否存在，如果目录不存在则创建
	if _, err := os.Stat(UploadFile); os.IsNotExist(err) {
		// 目录不存在，创建目录
		if err := os.MkdirAll(UploadFile, os.ModePerm); err != nil {
			msg = "文件夹创建失败"
			z.Error(fmt.Sprintf("文件夹创建失败：%s %s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		z.Info("Directory created successfully!")
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		msg = "文件保存失败"
		z.Error(fmt.Sprintf("保存失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	f, err := excelize.OpenFile(dst)
	if err != nil {
		msg = "文件打开失败"
		z.Error(fmt.Sprintf("文件打开失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			z.Error(fmt.Sprintf("关闭导入的文件失败：%s", err))
		}
		if err := os.Remove(dst); err != nil {
			z.Error(fmt.Sprintf("删除导入的文件失败：%s", err))
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		msg = "获取数据失败"
		z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 是否为 3 的倍数行
	if (len(rows))%3 != 0 {
		msg = "数据不合法"
		z.Error(fmt.Sprintf("数据不合法：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	var res []resBatchProcessItem

	round := 1
	for i := 0; i < len(rows); i = i + 3 {
		list, err := GetWordListByWordArray(rows[i])
		if err != nil {
			msg = "获取数据失败"
			z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}

		// 拼接 known 属性
		for j := 0; j < len(list); j++ {
			list[j].Known = rows[i+1][j]
		}

		// 此处开始估算词汇量
		vo := EstimateVocabulary(list)

		item := resBatchProcessItem{
			Round:            round,
			TestVocabulary:   vo,
			PreplyVocabulary: rows[i+2][0],
			CalculatedWords:  len(list),
		}
		res = append(res, item)

		round++
	}

	response.Success(c, res)
	return
}

func BatchProcess2(c *gin.Context) {
	var (
		msg string
		err error
	)

	file, err := c.FormFile("file")
	if err != nil {
		msg = "文件上传失败"
		z.Error(fmt.Sprintf("%s:%s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	//获取文件名
	fileName := file.Filename

	if file.Size > 10*1024*1024 {
		msg = "文件过大"
		z.Error(fmt.Sprintf("文件过大:%s", msg))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	newUUID, _ := uuid.NewUUID()

	// 文件后缀
	fileSuffix := path.Ext(fileName)
	// 路径
	dst := path.Join(UploadFile, newUUID.String()+fileSuffix)

	// 先判断目录是否存在，如果目录不存在则创建
	if _, err := os.Stat(UploadFile); os.IsNotExist(err) {
		// 目录不存在，创建目录
		if err := os.MkdirAll(UploadFile, os.ModePerm); err != nil {
			msg = "文件夹创建失败"
			z.Error(fmt.Sprintf("文件夹创建失败：%s %s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		z.Info("Directory created successfully!")
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		msg = "文件保存失败"
		z.Error(fmt.Sprintf("保存失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	f, err := excelize.OpenFile(dst)
	if err != nil {
		msg = "文件打开失败"
		z.Error(fmt.Sprintf("文件打开失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			z.Error(fmt.Sprintf("关闭导入的文件失败：%s", err))
		}
		if err := os.Remove(dst); err != nil {
			z.Error(fmt.Sprintf("删除导入的文件失败：%s", err))
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		msg = "获取数据失败"
		z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 是否为 3 的倍数行
	if (len(rows)+1)%3 != 0 {
		msg = "数据不合法"
		z.Error(fmt.Sprintf("数据不合法：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	var res []resBatchProcessItem

	round := 1
	for i := 0; i < len(rows); i = i + 3 {

		list, err := GetWordListByWordArray(rows[i])
		if err != nil {
			msg = "获取数据失败"
			z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}

		// 拼接 known 属性
		for j := 0; j < len(list); j++ {
			list[j].Known = rows[i+1][j]
		}

		// 此处开始估算词汇量
		vo := EstimateVocabulary(list)

		item := resBatchProcessItem{
			Round:           round,
			TestVocabulary:  vo,
			CalculatedWords: len(list),
		}
		res = append(res, item)

		round++
	}

	rule := map[int]string{
		0: "不认识比例10%，长度100",
		1: "不认识比例20%，长度100",
		2: "不认识比例30%，长度100",
		3: "不认识比例10%，长度200",
		4: "不认识比例20%，长度200",
		5: "不认识比例30%，长度200",
		6: "不认识比例10%，长度300",
		7: "不认识比例20%，长度300",
		8: "不认识比例30%，长度300",
	}

	for i := 0; i < len(res)/100; i = i + 1 {
		var data []float64
		for j := 0; j < 100; j++ {
			data = append(data, float64(res[i*100+j].TestVocabulary))
		}
		// 计算平均数
		mean := calculateMean(data)
		fmt.Printf("%v的平均数: %.2f\n", rule[i], mean)

		// 计算方差
		variance := calculateVariance(data, mean)
		fmt.Printf("%v的方差: %.2f\n", rule[i], variance)
	}

	// 写回文件
	start := 3
	for i := 0; i < len(res); i++ {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", start), res[i].Round)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", start), res[i].TestVocabulary)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", start), res[i].CalculatedWords)
		start = start + 3
	}

	// 保存修改后的Excel文件

	// 路径
	newUUID2, _ := uuid.NewUUID()
	outputPath := path.Join(UploadFile, newUUID2.String()+fileSuffix)
	if err := f.SaveAs(outputPath); err != nil {
		msg = "文件保存失败"
		z.Error(fmt.Sprintf("文件保存失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	// 返回生成的Excel文件
	c.File(outputPath)

	//response.Success(c, res)
	return
}

// 计算平均数
func calculateMean(data []float64) float64 {
	sum := 0.0
	n := len(data)

	for _, val := range data {
		sum += val
	}

	return sum / float64(n)
}

// 计算方差
func calculateVariance(data []float64, mean float64) float64 {
	sum := 0.0
	n := len(data)

	for _, val := range data {
		sum += math.Pow(val-mean, 2)
	}

	return sum / float64(n)
}

func EstimateVocabulary(list []Wordbook2) int {
	// 设置抽样参数
	numLayers := 10 // 词汇太少了，无法估算
	if len(list) < numLayers {
		return 20
	}
	//totalWords := 54000
	totalWords := int(list[len(list)-1].ID)

	// 定义层次权重，采用指数衰减的方式
	weights := []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.1, 0.1}
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}

	// 定义层次的得分权重
	knownWeights := []float64{2.0, 2.1, 2.2, 2.3, 2.35, 2.4, 2.45, 2.5, 2.6, 2.7}

	// 定义层次系数
	//levelCoefficient := []float64{0.95, 0.8, 0.85, 0.8, 0.75, 0.7}

	// 计算每个层次的分界线 = 划分层级
	boundaries := make([]int, numLayers+1)
	boundaries[0] = 1
	for i := 1; i <= numLayers; i++ {
		boundary := int(math.Round(float64(totalWords) * (weights[i-1] / totalWeight)))
		boundaries[i] = boundaries[i-1] + boundary
	}

	// 计算每个层次有多少单词认识 = 认识率
	knownPers := make([]float64, numLayers)
	totalPer := 0.0
	curId := 0
	length := len(list)

	for i := 0; i < numLayers; i++ {
		tmp := curId
		knownCount := 0
		layerBoundaryStart := boundaries[i]
		layerBoundaryEnd := boundaries[i+1]

		for ; curId < length && list[curId].ID < uint(layerBoundaryEnd); curId++ {
			if list[curId].ID > uint(layerBoundaryStart) && list[curId].Known == "1" {
				knownCount++
			}
		}

		// 计算该层级的认识率
		if curId > 0 {
			knownPer := float64(knownCount) / float64(curId-tmp)
			knownPers[i] = knownPer
		} else {
			knownPers[i] = 0.0
		}
	}

	// 避免某些中间层级没有单词，导致平均认识率加了NAN，为0，这里做一下处理
	skip := 0
	for _, num := range knownPers {
		if math.IsNaN(num) {
			skip++
			continue // 跳过NaN值
		}
		totalPer += num
	}

	totalPer = totalPer / float64(numLayers-skip)
	if totalPer == 0 {
		return 20
	}

	totalKnownWeights := 0.0
	for i := 0; i < len(knownPers); i++ {
		if math.IsNaN(knownPers[i]) {
			skip++
			continue // 跳过NaN值
		}
		totalKnownWeights += knownWeights[i]
	}

	// 根据认识率的得分权重计算每个层级的得分
	knownScores := make([]float64, numLayers)
	totalScore := 0.0
	for i := 0; i < numLayers; i++ {
		if math.IsNaN(knownPers[i]) {
			continue // 跳过NaN值
		}
		knownScores[i] = knownPers[i] * (knownWeights[i] / totalKnownWeights)
		totalScore += knownScores[i]
	}
	totalNum2 := math.Round(float64(totalWords) * totalScore)
	// float64 转换成 int
	totalNum := int(totalNum2)
	if totalNum > 54000 || totalNum < 0 {
		return 20
	}

	return totalNum

}

type resCoefficientFourGrade struct {
	Correlation  float64   `json:"correlation"`
	FourGrades   []float64 `json:"four_grades"`
	Vocabularies []float64 `json:"vocabularies"`
}

func GetCoefficientFourGrade(c *gin.Context) {
	var (
		msg string
		err error
	)

	users, err := GetFourGradeAndVocabulary()
	if err != nil {
		msg = "获取数据失败"
		z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	var fourGrades []float64
	var vocabularies []float64
	for _, user := range users {
		fourGrade, err := strconv.Atoi(user.FourGrade)
		if err != nil {
			continue
		}

		var array []int
		err = json.Unmarshal([]byte(user.BasicVocabulary), &array)
		if err != nil || len(array) == 0 {
			continue
		}

		//	取array的平均值
		var sum int
		for _, v := range array {
			sum += v
		}
		vocabulary := sum / len(array)

		fourGrades = append(fourGrades, float64(fourGrade))
		vocabularies = append(vocabularies, float64(vocabulary))

	}

	// 计算皮尔逊相关系数
	correlation := pearsonCorrelation(fourGrades, vocabularies)

	res := resCoefficientFourGrade{
		Correlation:  correlation,
		FourGrades:   fourGrades,
		Vocabularies: vocabularies,
	}

	response.Success(c, res)
	return
}

type resCoefficientSixGrade struct {
	Correlation  float64   `json:"correlation"`
	SixGrades    []float64 `json:"six_grades"`
	Vocabularies []float64 `json:"vocabularies"`
}

func GetCoefficientSixGrade(c *gin.Context) {
	var (
		msg string
		err error
	)

	users, err := GetSixGradeAndVocabulary()
	if err != nil {
		msg = "获取数据失败"
		z.Error(fmt.Sprintf("获取数据失败：%s %s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	var SixGrades []float64
	var vocabularies []float64
	for _, user := range users {
		fourGrade, err := strconv.Atoi(user.SixGrade)
		if err != nil {
			continue
		}

		var array []int
		err = json.Unmarshal([]byte(user.BasicVocabulary), &array)
		if err != nil || len(array) == 0 {
			continue
		}

		//	取array的平均值
		var sum int
		for _, v := range array {
			sum += v
		}
		vocabulary := sum / len(array)

		SixGrades = append(SixGrades, float64(fourGrade))
		vocabularies = append(vocabularies, float64(vocabulary))

	}

	// 计算皮尔逊相关系数
	correlation := pearsonCorrelation(SixGrades, vocabularies)

	res := resCoefficientSixGrade{
		Correlation:  correlation,
		SixGrades:    SixGrades,
		Vocabularies: vocabularies,
	}

	response.Success(c, res)
	return
}

// 计算皮尔逊相关系数
func pearsonCorrelation(x, y []float64) float64 {
	n := len(x)

	// 计算平均值
	xMean := mean(x)
	yMean := mean(y)

	// 计算偏差乘积的和
	deviationProductSum := 0.0
	for i := 0; i < n; i++ {
		deviationProductSum += (x[i] - xMean) * (y[i] - yMean)
	}

	// 计算偏差的平方和
	xDeviationSum := deviationSum(x, xMean)
	yDeviationSum := deviationSum(y, yMean)

	// 计算皮尔逊相关系数
	correlation := deviationProductSum / math.Sqrt(xDeviationSum*yDeviationSum)

	return correlation
}

// 计算平均值
func mean(data []float64) float64 {
	sum := 0.0
	n := len(data)
	for _, val := range data {
		sum += val
	}
	return sum / float64(n)
}

// 计算偏差的平方和
func deviationSum(data []float64, mean float64) float64 {
	sum := 0.0
	for _, val := range data {
		deviation := val - mean
		sum += deviation * deviation
	}
	return sum
}
