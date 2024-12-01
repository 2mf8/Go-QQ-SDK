package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/2mf8/Go-QQ-Client/dto"
	v1 "github.com/2mf8/Go-QQ-Client/openapi/v1"
	"github.com/2mf8/Go-QQ-Client/token"
	"github.com/2mf8/Go-QQ-Client/webhook"
)

func main() {
	webhook.InitLog()
	as := webhook.ReadSetting()
	v1.GetAccessToken(as.AppId, as.AppSecret)
	var ctx context.Context
	token := token.BotToken(as.AppId, as.Token, string(token.TypeBot))
	api := NewOpenAPI(token).WithTimeout(3 * time.Second)
	b, _ := json.Marshal(as)
	fmt.Println("配置", string(b))
	webhook.GroupAtMessageEventHandler = func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		fmt.Println(data.GroupId, data.Content)
		if strings.TrimSpace(data.Content) == "测试" {
			api.PostGroupMessage(ctx, data.GroupId, &dto.GroupMessageToCreate{
				Content: "成功",
				MsgID:   data.MsgId,
				MsgType: 0,
			})
		}
		return nil
	}
	webhook.C2CMessageEventHandler = func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		b, _ := json.Marshal(event)
		fmt.Println(string(b), data.Content)
		return nil
	}
	webhook.InitGin()
	select {}
}
