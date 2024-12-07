package dto

import "github.com/2mf8/Go-QQ-SDK/dto/keyboard"

// MessageToCreate 发送消息结构体定义
type MessageToCreate struct {
	Content string `json:"content,omitempty"`
	Embed   *Embed `json:"embed,omitempty"`
	Ark     *Ark   `json:"ark,omitempty"`
	Image   string `json:"image,omitempty"`
	// 要回复的消息id，为空是主动消息，公域机器人会异步审核，不为空是被动消息，公域机器人会校验语料
	MsgID            string                    `json:"msg_id,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	Markdown         *Markdown                 `json:"markdown,omitempty"`
	Keyboard         *keyboard.MessageKeyboard `json:"keyboard,omitempty"` // 消息按钮组件
	EventID          string                    `json:"event_id,omitempty"` // 要回复的事件id, 逻辑同MsgID
}

type C2CMsgType uint8

const (
	// C2CMsgTypeText 文本
	C2CMsgTypeText C2CMsgType = iota
	// C2CMsgTypeTextPic 图文混排
	C2CMsgTypeTextPic
	// C2CMsgTypeMarkdown markdown
	C2CMsgTypeMarkdown
	// C2CMsgTypeArk ark
	C2CMsgTypeArk
	// C2CMsgTypeEmbed embed
	C2CMsgTypeEmbed
	// C2CMsgTypeMedia media
	C2CMsgTypeMedia C2CMsgType = 7
)

// GroupMessageToCreate 消息类型： 0 是文本，1 图文混排，2  markdown， 3 ark，4 embed 7 富媒体
type GroupMessageToCreate = C2CMessageToCreate

// C2CMessageToCreate 创建人人、群消息的结构体定义
type C2CMessageToCreate struct {
	Content          string                    `json:"content,omitempty"`
	MsgType          C2CMsgType                `json:"msg_type"`
	Markdown         *Markdown                 `json:"markdown,omitempty"`
	Keyboard         *keyboard.MessageKeyboard `json:"keyboard,omitempty"` // 消息按钮组件
	Media            *FileInfo                 `json:"media,omitempty"`
	Ark              *Ark                      `json:"ark,omitempty"`
	Image            string                    `json:"image,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	EventID          EventType                 `json:"event_id,omitempty"` // 要回复的事件id, 逻辑同MsgID
	MsgID            string                    `json:"msg_id,omitempty"`
	MsgReq           uint                      `json:"msg_req,omitempty"`
}

type FileInfo struct {
	FileInfo string `json:"file_info,omitempty"`
}

// GroupRichMediaMessageToCreate 媒体类型：1 图片，2 视频，3 语音，4 文件（暂不开放） 资源格式要求： 图片：png/jpg，视频：mp4，语音：silk，
type GroupRichMediaMessageToCreate = C2CRichMediaMessageToCreate

// C2CRichMediaMessageToCreate 媒体类型：1 图片，2 视频，3 语音，4 文件（暂不开放） 资源格式要求： 图片：png/jpg，视频：mp4，语音：silk，
type C2CRichMediaMessageToCreate struct {
	FileType   int    `json:"file_type"`
	Url        string `json:"url"`
	SrvSendMsg bool   `json:"srv_send_msg"`
	FileData   []byte `json:"file_data"`
}

// MessageReference 引用消息
type MessageReference struct {
	MessageID             string `json:"message_id"`               // 消息 id
	IgnoreGetMessageError bool   `json:"ignore_get_message_error"` // 是否忽律获取消息失败错误
}

// Markdown markdown 消息
type Markdown struct {
	TemplateID       int               `json:"template_id,omitempty"` // 模版 id
	CustomTemplateId string            `json:"custom_template_id,omitempty"`
	Params           []*MarkdownParams `json:"params"`  // 模版参数
	Content          string            `json:"content"` // 原生 markdown
}

// MarkdownParams markdown 模版参数 键值对
type MarkdownParams struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// SettingGuideToCreate 发送引导消息的结构体
type SettingGuideToCreate struct {
	Content      string        `json:"content,omitempty"`       // 频道内发引导消息可以带@
	SettingGuide *SettingGuide `json:"setting_guide,omitempty"` // 设置引导
}

// SettingGuide 设置引导
type SettingGuide struct {
	// 频道ID, 当通过私信发送设置引导消息时，需要指定guild_id
	GuildID string `json:"guild_id"`
}

type RichMediaMsgResp struct {
	FileUuid string `json:"file_uuid,omitempty"`
	FileInfo string `json:"file_info,omitempty"`
	Ttl      uint   `json:"ttl,omitempty"`
}

type GroupMsgResp = C2CMsgResp

type C2CMsgResp struct {
	Id        string    `json:"id"`
	Timestamp Timestamp `json:"timestamp"`
}
