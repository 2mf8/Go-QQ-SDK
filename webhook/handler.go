package webhook

import "github.com/2mf8/Go-QQ-Client/dto"

// ReadyHandler 可以处理 ws 的 ready 事件
var ReadyHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSReadyData)

// ErrorNotifyHandler 当 ws 连接发生错误的时候，会回调，方便使用方监控相关错误
// 比如 reconnect invalidSession 等错误，错误可以转换为 bot.Err
var ErrorNotifyHandler func(err error)

// PlainEventHandler 透传handler
var PlainEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, message []byte) error

// GuildEventHandler 频道事件handler
var GuildEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGuildData) error

// GuildMemberEventHandler 频道成员事件 handler
var GuildMemberEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGuildMemberData) error

// ChannelEventHandler 子频道事件 handler
var ChannelEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSChannelData) error

// CheckEventHandler 消息前置检测
var CheckEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, message []byte) bool

// MessageEventHandler 消息事件 handler
var MessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageData) error

// MessageDeleteEventHandler 消息事件 handler
var MessageDeleteEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageDeleteData) error

// PublicMessageDeleteEventHandler 消息事件 handler
var PublicMessageDeleteEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSPublicMessageDeleteData) error

// DirectMessageDeleteEventHandler 消息事件 handler
var DirectMessageDeleteEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSDirectMessageDeleteData) error

// MessageReactionEventHandler 表情表态事件 handler
var MessageReactionEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageReactionData) error

// ATMessageEventHandler at 机器人消息事件 handler
var ATMessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSATMessageData) error

// DirectMessageEventHandler 私信消息事件 handler
var DirectMessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSDirectMessageData) error

// AudioEventHandler 音频机器人事件 handler
var AudioEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSAudioData) error

// MessageAuditEventHandler 消息审核事件 handler
var MessageAuditEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageAuditData) error

// ThreadEventHandler 论坛主题事件 handler
var ThreadEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSThreadData) error

// PostEventHandler 论坛回帖事件 handler
var PostEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSPostData) error

// ReplyEventHandler 论坛帖子回复事件 handler
var ReplyEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSReplyData) error

// ForumAuditEventHandler 论坛帖子审核事件 handler
var ForumAuditEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSForumAuditData) error

// InteractionEventHandler 互动事件 handler
var InteractionEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSInteractionData) error

var GroupAtMessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupATMessageData) error

var GroupMessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMessageData) error

var C2CMessageEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSC2CMessageData) error

var GroupAddRobotEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupAddRobotData) error

var GroupDelRobotEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupDelRobotData) error

var GroupMsgRejectEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMsgRejectData) error

var GroupMsgReceiveEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMsgReceiveData) error

var FriendAddEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendAddData) error

var FriendDelEventHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendDelData) error

var C2CMsgRejectHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendMsgRejectData) error

var C2CMsgReceiveHandler func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendMsgReveiceData) error

func init() {
	// ReadyHandler 可以处理 ws 的 ready 事件
	ReadyHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSReadyData) {}

	// ErrorNotifyHandler 当 ws 连接发生错误的时候，会回调，方便使用方监控相关错误
	// 比如 reconnect invalidSession 等错误，错误可以转换为 bot.Err
	ErrorNotifyHandler = func(err error) {}

	// PlainEventHandler 透传handler
	PlainEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, message []byte) error {
		return nil
	}

	// GuildEventHandler 频道事件handler
	GuildEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGuildData) error {
		return nil
	}

	// GuildMemberEventHandler 频道成员事件 handler
	GuildMemberEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		return nil
	}

	// ChannelEventHandler 子频道事件 handler
	ChannelEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSChannelData) error {
		return nil
	}

	// CheckEventHandler 消息前置检测
	CheckEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, message []byte) bool {
		return false
	}

	// MessageEventHandler 消息事件 handler
	MessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageData) error {
		return nil
	}

	// MessageDeleteEventHandler 消息事件 handler
	MessageDeleteEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageDeleteData) error {
		return nil
	}

	// PublicMessageDeleteEventHandler 消息事件 handler
	PublicMessageDeleteEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSPublicMessageDeleteData) error {
		return nil
	}

	// DirectMessageDeleteEventHandler 消息事件 handler
	DirectMessageDeleteEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSDirectMessageDeleteData) error {
		return nil
	}

	// MessageReactionEventHandler 表情表态事件 handler
	MessageReactionEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageReactionData) error {
		return nil
	}

	// ATMessageEventHandler at 机器人消息事件 handler
	ATMessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSATMessageData) error {
		return nil
	}

	// DirectMessageEventHandler 私信消息事件 handler
	DirectMessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		return nil
	}

	// AudioEventHandler 音频机器人事件 handler
	AudioEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSAudioData) error {
		return nil
	}

	// MessageAuditEventHandler 消息审核事件 handler
	MessageAuditEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSMessageAuditData) error {
		return nil
	}

	// ThreadEventHandler 论坛主题事件 handler
	ThreadEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSThreadData) error {
		return nil
	}

	// PostEventHandler 论坛回帖事件 handler
	PostEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSPostData) error {
		return nil
	}

	// ReplyEventHandler 论坛帖子回复事件 handler
	ReplyEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSReplyData) error {
		return nil
	}

	// ForumAuditEventHandler 论坛帖子审核事件 handler
	ForumAuditEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSForumAuditData) error {
		return nil
	}

	// InteractionEventHandler 互动事件 handler
	InteractionEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSInteractionData) error {
		return nil
	}

	GroupAtMessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		return nil
	}

	GroupMessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMessageData) error {
		return nil
	}

	C2CMessageEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		return nil
	}

	GroupAddRobotEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupAddRobotData) error {
		return nil
	}

	GroupDelRobotEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupDelRobotData) error {
		return nil
	}

	GroupMsgRejectEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMsgRejectData) error {
		return nil
	}

	GroupMsgReceiveEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSGroupMsgReceiveData) error {
		return nil
	}

	FriendAddEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendAddData) error {
		return nil
	}

	FriendDelEventHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendDelData) error {
		return nil
	}

	C2CMsgRejectHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendMsgRejectData) error {
		return nil
	}

	C2CMsgReceiveHandler = func(bot *BotHeaderInfo, event *dto.WSPayload, data *dto.WSFriendMsgReveiceData) error {
		return nil
	}
}
