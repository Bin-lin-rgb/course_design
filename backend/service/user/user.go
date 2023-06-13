package user

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var (
	response *utils.Response
)

type reqRegister struct {
	Account  string `form:"account" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var (
		reqForm reqRegister
		msg     string
		err     error
	)
	if err = c.ShouldBind(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	user := UserInfo{
		Account:  reqForm.Account,
		Username: reqForm.Username,
		Password: reqForm.Password,
	}

	err = user.Create()
	if err != nil {
		if utils.IsUniqueConstraintError(err) {
			msg = "用户名重复"
			z.Error(msg)
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		msg = "数据库插入错误"
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	response.Success(c, nil)
}

type reqLogin struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var (
		reqForm reqLogin
		msg     string
		err     error
	)
	if err = c.ShouldBind(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	user := UserInfo{
		Account:  reqForm.Account,
		Password: reqForm.Password,
	}
	_, found := Cache.Get(LoginErrCountKey + reqForm.Account)

	err = user.IsExist()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//如果能找到，就增加一次失败次数
			if found == true {
				_, _ = Cache.IncrementInt(LoginErrCountKey+reqForm.Account, 1)

			} else {
				//没有记录，说明是第一次登录
				Cache.Set(LoginErrCountKey+reqForm.Account, 1, 5*time.Minute)
			}

			errCount, _ := Cache.Get(LoginErrCountKey + reqForm.Account)
			//重新设置5分钟过期时间
			if errCount == MaxLoginErrCount {
				Cache.Set(LoginErrCountKey+reqForm.Account, MaxLoginErrCount, 5*time.Minute)
			}
			if errCount.(int) > MaxLoginErrCount {
				msg = "密码错误 3 次，请 5 分钟后再试"
				z.Error(fmt.Sprintf("%s", msg))
				response.Err(c, http.StatusOK, msg, nil)
				return
			}

			msg = "不存在该用户"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf("%s:%s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	token, err := j.TokenNew(int64(user.ID), user.Username) //新建token
	if err != nil {
		msg = "新建 token 失败"
		z.Error(fmt.Sprintf("%s:%s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	resp := struct {
		Token string `json:"token,omitempty"`
	}{Token: token}
	response.Success(c, resp)
}

type reqSetGrade struct {
	FourGrade string `json:"fourGrade" binding:"required"`
	SixGrade  string `json:"sixGrade" binding:"required"`
}

func SetGrade(c *gin.Context) {
	var (
		reqForm reqSetGrade
		msg     string
		err     error
	)

	userID, _ := c.Get("userID")

	if err = c.ShouldBind(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	user := UserInfo{
		FourGrade: reqForm.FourGrade,
		SixGrade:  reqForm.SixGrade,
		Model:     gorm.Model{ID: uint(userID.(int64))},
	}
	err = user.SetGrade()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg = fmt.Sprintf("账户 %s 不存在，无法设置四级六级成绩", user.Account)
			z.Error(msg)
		} else {
			msg = "设置四六级成绩失败"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
		}
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	response.Success(c, nil)
}

type reqResetPassword struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	var (
		reqForm reqResetPassword
		msg     string
		err     error
	)
	if err = c.ShouldBind(&reqForm); err != nil {
		msg = "请求不合法"
		z.Error(msg)
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	user := UserInfo{
		Account:  reqForm.Account,
		Password: reqForm.Password,
	}
	err = user.UpdatePassword()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg = fmt.Sprintf("账户 %s 不存在，无法修改密码", user.Account)
			z.Error(msg)
		} else {
			msg = "更新密码失败"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
		}
		response.Err(c, http.StatusOK, msg, nil)
		return
	}
	response.Success(c, nil)
}

func GetUserInfo(c *gin.Context) {
	var (
		msg string
		err error
	)
	userID, _ := c.Get("userID")

	user := UserInfo{
		Model: gorm.Model{ID: uint(userID.(int64))},
	}
	err = user.GetUserInfoByID()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg = "不存在该用户"
			z.Error(fmt.Sprintf("%s:%s", msg, err))
			response.Err(c, http.StatusOK, msg, nil)
			return
		}
		msg = "数据库查询失败"
		z.Error(fmt.Sprintf("%s:%s", msg, err))
		response.Err(c, http.StatusOK, msg, nil)
		return
	}

	resp := struct {
		UserId   int64  `json:"userId,omitempty"`
		Username string `json:"username,omitempty"`
	}{
		UserId:   userID.(int64),
		Username: user.Username,
	}

	response.Success(c, resp)
}
