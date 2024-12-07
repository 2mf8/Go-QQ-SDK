# Go-QQ-SDK

QQ频道机器人，官方 GOLANG SDK。

[![Go Reference](https://pkg.go.dev/badge/github.com/2mf8/Go-QQ-SDK.svg)](https://pkg.go.dev/github.com/2mf8/Go-QQ-SDK)

# [QQ交流群 677742758](https://qm.qq.com/q/okWktIaAqk)

<details>

<summary><font size="4">已完成功能/开发计划列表</font></summary>

### **登录**

- [x] 登录

### **消息类型**
- [x] 文本
- [x] 图片
- [x] 语音
- [x] MarkDown
- [ ] 表情
- [ ] At
- [ ] 回复
- [ ] 长消息(仅群聊/私聊)
- [ ] 链接分享
- [ ] 小程序(暂只支持RAW)
- [x] 短视频
- [ ] 合并转发
- [ ] 群文件(上传与接收信息)

### **群聊**

- [x] 收发群消息
- [x] 机器人加群通知
- [x] 机器人离群通知
- [x] 群接收机器人主动消息通知
- [x] 群拒绝机器人主动消息通知

### **C2C**

- [x] 收发C2C消息
- [x] 机器人加好友通知
- [x] 机器人删好友通知
- [x] 接收机器人消息通知
- [x] 拒绝机器人消息通知

</details>

## 一、如何使用

### 1.回调地址配置

https://你的域名:端口/qqbot/你的应用appid

示例 `https://fw1009zb5979.vicp.fun:443/qqbot/101981675`

### 2.配置文件填写（支持多账号）

默认配置文件为
```
{
	"apps": {
		"123456": {
			"qq": 123456,
			"app_id": 123456,
			"token": "你的AppToken",
			"app_secret": "你的AppSecret"
		}
	},
	"port": 8443,
	"cert_file": "ssl证书文件路径",
	"cert_key": "ssl证书密钥"
}
```
多账号

```
{
	"apps": {
		"5123456": {
			"qq": 123456,
			"app_id": 5123456,
			"token": "你的AppToken",
			"app_secret": "你的AppSecret"
		},
		"7234567": {
			"qq": 234567,
			"app_id": 7234567,
			"token": "你的AppToken",
			"app_secret": "你的AppSecret"
		}
	},
	"port": 8443,
	"cert_file": "ssl证书文件路径",
	"cert_key": "ssl证书密钥"
}
```

### 3.请求 openapi 接口，操作资源

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/2mf8/Go-QQ-SDK/dto"
	"github.com/2mf8/Go-QQ-SDK/openapi"
	"github.com/2mf8/Go-QQ-SDK/token"
	"github.com/2mf8/Go-QQ-SDK/webhook"
	log "github.com/sirupsen/logrus"
)

var Apis = make(map[string]openapi.OpenAPI, 0)

func main() {
	webhook.InitLog()
	as := webhook.ReadSetting()
	var ctx context.Context
	for i, v := range as.Apps {
		token := token.BotToken(v.AppId, v.Token, string(token.TypeBot))
		api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
		Apis[i] = api
	}
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	webhook.GroupAtMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		fmt.Println(bot.XBotAppid, data.GroupId, data.Content)
		if len(data.Attachments) > 0 {
			log.Infof(`BotId(%s) GroupId(%s) UserId(%s) <- %s <image id="%s">`, bot.XBotAppid[0], data.GroupId, data.Author.UserId, data.Content, data.Attachments[0].URL)
		} else {
			log.Infof("BotId(%s) GroupId(%s) UserId(%s) <- %s", bot.XBotAppid[0], data.GroupId, data.Author.UserId, data.Content)
		}
		if strings.TrimSpace(data.Content) == "测试" {
			Apis[bot.XBotAppid[0]].PostGroupMessage(ctx, data.GroupId, &dto.GroupMessageToCreate{
				Content: "成功",
				MsgID:   data.MsgId,
				MsgType: 0,
			})
		}
		return nil
	}
	webhook.C2CMessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		b, _ := json.Marshal(event)
		fmt.Println(bot.XBotAppid, string(b), data.Content)
		return nil
	}
	webhook.MessageEventHandler = func(bot *webhook.BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageData) error {
		b, _ := json.Marshal(event)
		fmt.Println(bot.XBotAppid, string(b), data.Content)
		return nil
	}
	webhook.InitGin()
	select {}
}
```