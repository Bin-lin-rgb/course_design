package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	vp           *viper.Viper
	GlobalConfig GlobalConf
)

type WebSeverConf struct {
	HttpsListenPort string `json:"httpsListenPort,omitempty"`
}

type MysqlDbConf struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	DbName   string `json:"dbName,omitempty"`
}

type GlobalConf struct {
	WebSever WebSeverConf `json:"webSever"`
	MysqlDb  MysqlDbConf  `json:"mysqlDb"`
}

func init() {
	//初始化viper
	vp = viper.New()
	//设置文件名
	vp.SetConfigName("config")

	//设置文件类型
	vp.SetConfigType("json")
	//设置文件所在的目录
	vp.AddConfigPath("./config")

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；
			panic("config file not found")
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic("init GlobalConfig error")
		}
	}
	if err := vp.Unmarshal(&GlobalConfig); err != nil {
		panic("读取配置文件出错")
	}
}

func GetConfig() {
	// 将读取的配置信息保存至全局变量GlobalConfig
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic(fmt.Errorf("unmarshal GlobalConfig failed, err:%s \n", err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	//配置文件发生变化后要同步到全局变量config
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("GlobalConfig file has been modify")
		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})
}
