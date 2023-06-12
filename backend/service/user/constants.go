package user

import (
	"backend/common"
	"backend/utils"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	z  = common.GetLogger()
	db = common.GetDbConn()
	j  *utils.JWT

	MaxLoginErrCount = 3
	LoginErrCountKey = "LOGIN_ERR_COUNT_KEY"

	Cache = cache.New(5*time.Minute, 10*time.Minute)
)
