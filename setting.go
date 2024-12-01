package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/2mf8/Go-QQ-Client/log"
)

type Setting struct {
	QQ        uint64
	AppId     uint64
	Token     string
	AppSecret string
}

var SettingPath = "setting"
var AllSetting *Setting = &Setting{}
var NullSetting = &Setting{}

func AllSettings() *Setting {
	_, err := toml.DecodeFile("setting/setting.toml", AllSetting)
	if err != nil {
		return AllSetting
	}
	return AllSetting
}

func ReadSetting() Setting {
	tomlData := `QQ = 0
AppId = 0
Token = ""
AppSecret = ""
	`
	if !PathExists(SettingPath) {
		if err := os.MkdirAll(SettingPath, 0777); err != nil {
			log.Warnf("failed to mkdir")
			return *AllSetting
		}
	}
	_, err := os.Stat(fmt.Sprintf("%s/setting.toml", SettingPath))
	if err != nil {
		err = os.WriteFile(fmt.Sprintf("%s/setting.toml", SettingPath), []byte(tomlData), 0644)
		if err != nil {
			log.Warn(err)
			log.Warn("生成配置文件 conf.toml 失败,请稍后重试。")
			return *AllSetting
		}
		log.Warn("已生成配置文件 conf.toml ,请修改后重新启动程序。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	AllSetting = AllSettings()
	return *AllSetting
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}