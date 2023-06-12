package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"net/http"
)

type respCaptcha struct {
	Id         string `json:"id"`
	Base64Blob string `json:"base64Blob"`
}

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成图形化验证码
func GenerateCaptcha(c *gin.Context) {
	var (
		msg string
		err error
	)
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          80,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	// 自定义配置
	driverString = captchaConfig
	driver = driverString.ConvertFonts()

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		msg = fmt.Sprintf("生成验证码失败: %v", err)
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
	}

	resp := respCaptcha{
		Id:         id,
		Base64Blob: b64s,
	}

	response.Success(c, resp)
}

type reqVerifyCode struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

// VerifyCode 验证captcha是否正确
func VerifyCode(c *gin.Context) {
	var (
		reqForm reqVerifyCode
		msg     string
		err     error
	)

	if err = c.ShouldBind(&reqForm); err != nil {
		msg = fmt.Sprintf("请求不合法: %v", err)
		z.Error(msg)
		response.Err(c, http.StatusOK, "请求不合法", nil)
		return
	}

	if !store.Verify(reqForm.Id, reqForm.Code, true) {
		msg = "验证码错误，请刷新验证码并重新输入"
		z.Error(fmt.Sprintf("验证码错误: %v", reqForm))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	response.Success(c, nil)
}
