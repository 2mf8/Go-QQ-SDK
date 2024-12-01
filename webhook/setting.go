package webhook

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type Setting struct {
	QQ        uint64 `json:"qq,omitempty" toml:"QQ"`
	AppId     uint64 `json:"app_id,omitempty" toml:"AppId"`
	Token     string `json:"token,omitempty" toml:"Token"`
	AppSecret string `json:"app_secret,omitempty" toml:"AppSecret"`
	Port      int    `json:"port,omitempty" toml:"Port"`
	CertFile  string `json:"cert_file,omitempty" toml:"CertFile"`
	CertKey   string `json:"cert_key,omitempty" toml:"CertKey"`
}

var SettingPath = "setting"
var AllSetting *Setting = &Setting{}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

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
Token = "你的AppToken"
AppSecret = "你的AppSecret"
Port = 8443
CertFile = "ssl证书文件路径"
CertKey = "ssl证书密钥"
	`
	if !PathExists(SettingPath) {
		if err := os.MkdirAll(SettingPath, 0777); err != nil {
			log.Warnf("failed to mkdir")
			return *AllSetting
		}
	}
	_, err := os.Stat(fmt.Sprintf("%s/setting.toml", SettingPath))
	if err != nil {
		_ = os.WriteFile(fmt.Sprintf("%s/setting.toml", SettingPath), []byte(tomlData), 0644)
		log.Warn("已生成配置文件 conf.toml 。")
		log.Info("请修改 conf.toml 后重新启动。")
		time.Sleep(time.Second * 5)
		//os.Exit(1)
	}
	AllSetting = AllSettings()
	return *AllSetting
}
