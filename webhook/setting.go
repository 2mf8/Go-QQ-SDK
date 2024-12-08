package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type App struct {
	QQ        uint64 `json:"qq,omitempty" toml:"QQ"`
	AppId     uint64 `json:"app_id,omitempty" toml:"AppId"`
	Token     string `json:"token,omitempty" toml:"Token"`
	AppSecret string `json:"app_secret,omitempty" toml:"AppSecret"`
}
type Setting struct {
	Apps     map[string]*App `json:"apps,omitempty" toml:"Apps"`
	Port     int             `json:"port,omitempty" toml:"Port"`
	CertFile string          `json:"cert_file,omitempty" toml:"CertFile"`
	CertKey  string          `json:"cert_key,omitempty" toml:"CertKey"`
}

var SettingPath = "setting"
var AllSetting *Setting = &Setting{}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func AllSettings() *Setting {
	s := &Setting{}
	jsonFile, err := os.Open(fmt.Sprintf("%s/setting.json", SettingPath))
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return s
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return s
	}
	json.Unmarshal(jsonData, &s)
	return s
}

func ReadSetting() Setting {
	app := &App{
		QQ:        123456,
		AppId:     123456,
		Token:     "你的AppToken",
		AppSecret: "你的AppSecret",
	}
	appMap := make(map[string]*App)
	appMap[fmt.Sprintf("%v", app.AppId)] = app
	apps := &Setting{
		Apps:     appMap,
		Port:     8443,
		CertFile: "ssl证书文件路径",
		CertKey:  "ssl证书密钥",
	}

	output, err := json.MarshalIndent(apps, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return *apps
	}

	if !PathExists(SettingPath) {
		if err := os.MkdirAll(SettingPath, 0777); err != nil {
			log.Warnf("failed to mkdir")
			return *AllSetting
		}
	}
	_, err = os.Stat(fmt.Sprintf("%s/setting.json", SettingPath))
	if err != nil {
		_ = os.WriteFile(fmt.Sprintf("%s/setting.json", SettingPath), []byte(output), 0644)
		log.Warn("已生成配置文件 setting.json 。")
		log.Info("请修改 setting.json 后重新启动。")
		log.Info("程序 10 秒后退出")
		time.Sleep(time.Second * 10)
		os.Exit(1)
	}
	AllSetting = AllSettings()
	return *AllSetting
}
