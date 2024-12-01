package webhook

import "github.com/2mf8/Go-QQ-Client/dto"

// ReadyHandler 可以处理 ws 的 ready 事件
var ReadyHandler func(event *dto.WSPayload, data *dto.WSReadyData)

// ErrorNotifyHandler 当 ws 连接发生错误的时候，会回调，方便使用方监控相关错误
// 比如 reconnect invalidSession 等错误，错误可以转换为 bot.Err
var ErrorNotifyHandler func(err error)

// PlainEventHandler 透传handler
var PlainEventHandler func(event *dto.WSPayload, message []byte) error

// GuildEventHandler 频道事件handler
var GuildEventHandler func(event *dto.WSPayload, data *dto.WSGuildData) error

// GuildMemberEventHandler 频道成员事件 handler
var GuildMemberEventHandler func(event *dto.WSPayload, data *dto.WSGuildMemberData) error

// ChannelEventHandler 子频道事件 handler
var ChannelEventHandler func(event *dto.WSPayload, data *dto.WSChannelData) error

// CheckEventHandler 消息前置检测
var CheckEventHandler func(event *dto.WSPayload, message []byte) bool

// MessageEventHandler 消息事件 handler
var MessageEventHandler func(event *dto.WSPayload, data *dto.WSMessageData) error

// MessageDeleteEventHandler 消息事件 handler
var MessageDeleteEventHandler func(event *dto.WSPayload, data *dto.WSMessageDeleteData) error

// PublicMessageDeleteEventHandler 消息事件 handler
var PublicMessageDeleteEventHandler func(event *dto.WSPayload, data *dto.WSPublicMessageDeleteData) error

// DirectMessageDeleteEventHandler 消息事件 handler
var DirectMessageDeleteEventHandler func(event *dto.WSPayload, data *dto.WSDirectMessageDeleteData) error

// MessageReactionEventHandler 表情表态事件 handler
var MessageReactionEventHandler func(event *dto.WSPayload, data *dto.WSMessageReactionData) error

// ATMessageEventHandler at 机器人消息事件 handler
var ATMessageEventHandler func(event *dto.WSPayload, data *dto.WSATMessageData) error

// DirectMessageEventHandler 私信消息事件 handler
var DirectMessageEventHandler func(event *dto.WSPayload, data *dto.WSDirectMessageData) error

// AudioEventHandler 音频机器人事件 handler
var AudioEventHandler func(event *dto.WSPayload, data *dto.WSAudioData) error

// MessageAuditEventHandler 消息审核事件 handler
var MessageAuditEventHandler func(event *dto.WSPayload, data *dto.WSMessageAuditData) error

// ThreadEventHandler 论坛主题事件 handler
var ThreadEventHandler func(event *dto.WSPayload, data *dto.WSThreadData) error

// PostEventHandler 论坛回帖事件 handler
var PostEventHandler func(event *dto.WSPayload, data *dto.WSPostData) error

// ReplyEventHandler 论坛帖子回复事件 handler
var ReplyEventHandler func(event *dto.WSPayload, data *dto.WSReplyData) error

// ForumAuditEventHandler 论坛帖子审核事件 handler
var ForumAuditEventHandler func(event *dto.WSPayload, data *dto.WSForumAuditData) error

// InteractionEventHandler 互动事件 handler
var InteractionEventHandler func(event *dto.WSPayload, data *dto.WSInteractionData) error

var GroupAtMessageEventHandler func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error

var GroupMessageEventHandler func(event *dto.WSPayload, data *dto.WSGroupMessageData) error

var C2CMessageEventHandler func(event *dto.WSPayload, data *dto.WSC2CMessageData) error

var GroupAddRobotEventHandle func(event *dto.WSPayload, data *dto.WSGroupAddRobotData) error

var GroupDelRobotEventHandle func(event *dto.WSPayload, data *dto.WSGroupDelRobotData) error

var GroupMsgRejectEventHandle func(event *dto.WSPayload, data *dto.WSGroupMsgRejectData) error

var GroupMsgReceiveEventHandle func(event *dto.WSPayload, data *dto.WSGroupMsgReceiveData) error

var FriendAddEventHandle func(event *dto.WSPayload, data *dto.WSFriendAddData) error

var FriendDelEventHandle func(event *dto.WSPayload, data *dto.WSFriendDelData) error

var C2CMsgRejectHandle func(event *dto.WSPayload, data *dto.WSFriendMsgRejectData) error

var C2CMsgReceiveHandle func(event *dto.WSPayload, data *dto.WSFriendMsgReveiceData) error
