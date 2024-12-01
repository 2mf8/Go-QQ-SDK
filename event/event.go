package event

import (
	"encoding/json"

	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/tidwall/gjson" // 由于回包的 d 类型不确定，gjson 用于从回包json中提取 d 并进行针对性的解析
)

var eventParseFuncMap = map[dto.OPCode]map[dto.EventType]eventParseFunc{
	dto.WSDispatchEvent: {
		dto.EventGuildCreate: guildHandler,
		dto.EventGuildUpdate: guildHandler,
		dto.EventGuildDelete: guildHandler,

		dto.EventChannelCreate: channelHandler,
		dto.EventChannelUpdate: channelHandler,
		dto.EventChannelDelete: channelHandler,

		dto.EventGuildMemberAdd:    guildMemberHandler,
		dto.EventGuildMemberUpdate: guildMemberHandler,
		dto.EventGuildMemberRemove: guildMemberHandler,

		dto.EventMessageCreate: messageHandler,
		dto.EventMessageDelete: messageDeleteHandler,

		dto.EventMessageReactionAdd:    messageReactionHandler,
		dto.EventMessageReactionRemove: messageReactionHandler,

		dto.EventAtMessageCreate:     atMessageHandler,
		dto.EventPublicMessageDelete: publicMessageDeleteHandler,

		dto.EventDirectMessageCreate: directMessageHandler,
		dto.EventDirectMessageDelete: directMessageDeleteHandler,

		dto.EventAudioStart:  audioHandler,
		dto.EventAudioFinish: audioHandler,
		dto.EventAudioOnMic:  audioHandler,
		dto.EventAudioOffMic: audioHandler,

		dto.EventMessageAuditPass:   messageAuditHandler,
		dto.EventMessageAuditReject: messageAuditHandler,

		dto.EventForumThreadCreate: threadHandler,
		dto.EventForumThreadUpdate: threadHandler,
		dto.EventForumThreadDelete: threadHandler,
		dto.EventForumPostCreate:   postHandler,
		dto.EventForumPostDelete:   postHandler,
		dto.EventForumReplyCreate:  replyHandler,
		dto.EventForumReplyDelete:  replyHandler,
		dto.EventForumAuditResult:  forumAuditHandler,

		dto.EventInteractionCreate: interactionHandler,

		dto.EventGroupATMessageCreate: groupAtMessageHandler,
		dto.EventGroupMessageCreate:   groupMessageHandler,
		dto.EventC2CMessageCreate:     c2cMessageHandle,

		dto.EventGroupAddRobbot:  groupAddRobotHandle,
		dto.EventGroupDelRobbot:  groupDelRobotHandle,
		dto.EventGroupMsgReject:  groupMsgRejectHandle,
		dto.EventGroupMsgReceive: groupMsgReceiveHandle,

		dto.EventFriendAdd:     friendAddHandle,
		dto.EventFriendDel:     friendDelHandle,
		dto.EventC2CMsgReject:  c2cMsgRejectHandle,
		dto.EventC2CMsgReceive: c2cMsgReceiveHandle,
	},
}

type eventParseFunc func(event *dto.WSPayload, message []byte) error

// ParseAndHandle 处理回调事件
func ParseAndHandle(payload *dto.WSPayload) error {
	// 指定类型的 handler
	if (DefaultHandlers.Check != nil) && DefaultHandlers.Check(payload, payload.RawMessage) == false {
		return nil
	}
	if h, ok := eventParseFuncMap[payload.OPCode][payload.Type]; ok {
		return h(payload, payload.RawMessage)
	}
	// 透传handler，如果未注册具体类型的 handler，会统一投递到这个 handler
	if DefaultHandlers.Plain != nil {
		return DefaultHandlers.Plain(payload, payload.RawMessage)
	}
	return nil
}

// ParseData 解析数据
func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	err := json.Unmarshal([]byte(data.String()), target)
	return err
}

// ParseData 解析数据
func GetId(message []byte) string {
	data := gjson.Get(string(message), "id")
	return data.String()
}

func guildHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Guild != nil {
		return DefaultHandlers.Guild(payload, data)
	}
	return nil
}

func channelHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSChannelData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Channel != nil {
		return DefaultHandlers.Channel(payload, data)
	}
	return nil
}

func guildMemberHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildMemberData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GuildMember != nil {
		return DefaultHandlers.GuildMember(payload, data)
	}
	return nil
}

func messageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Message != nil {
		return DefaultHandlers.Message(payload, data)
	}
	return nil
}

func messageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.MessageDelete != nil {
		return DefaultHandlers.MessageDelete(payload, data)
	}
	return nil
}

func messageReactionHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageReactionData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.MessageReaction != nil {
		return DefaultHandlers.MessageReaction(payload, data)
	}
	return nil
}

func atMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSATMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.ATMessage != nil {
		return DefaultHandlers.ATMessage(payload, data)
	}
	return nil
}

func publicMessageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSPublicMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.PublicMessageDelete != nil {
		return DefaultHandlers.PublicMessageDelete(payload, data)
	}
	return nil
}

func directMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSDirectMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.DirectMessage != nil {
		return DefaultHandlers.DirectMessage(payload, data)
	}
	return nil
}

func directMessageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSDirectMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.DirectMessageDelete != nil {
		return DefaultHandlers.DirectMessageDelete(payload, data)
	}
	return nil
}

func audioHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSAudioData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Audio != nil {
		return DefaultHandlers.Audio(payload, data)
	}
	return nil
}

func threadHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSThreadData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Thread != nil {
		return DefaultHandlers.Thread(payload, data)
	}
	return nil
}

func postHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSPostData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Post != nil {
		return DefaultHandlers.Post(payload, data)
	}
	return nil
}

func replyHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSReplyData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.Reply != nil {
		return DefaultHandlers.Reply(payload, data)
	}
	return nil
}

func forumAuditHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSForumAuditData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.ForumAudit != nil {
		return DefaultHandlers.ForumAudit(payload, data)
	}
	return nil
}

func messageAuditHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageAuditData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.MessageAudit != nil {
		return DefaultHandlers.MessageAudit(payload, data)
	}
	return nil
}

func interactionHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSInteractionData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	data.ID = GetId(message)
	if DefaultHandlers.Interaction != nil {
		return DefaultHandlers.Interaction(payload, data)
	}
	return nil
}

func groupAtMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupATMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupAtMessage != nil {
		return DefaultHandlers.GroupAtMessage(payload, data)
	}
	return nil
}

func groupMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupMessage != nil {
		return DefaultHandlers.GroupMessage(payload, data)
	}
	return nil
}

func c2cMessageHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSC2CMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.C2CMessage != nil {
		return DefaultHandlers.C2CMessage(payload, data)
	}
	return nil
}

func groupAddRobotHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupAddRobotData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupAddRobot != nil {
		return DefaultHandlers.GroupAddRobot(payload, data)
	}
	return nil
}

func groupDelRobotHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupDelRobotData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupDelRobot != nil {
		return DefaultHandlers.GroupDelRobot(payload, data)
	}
	return nil
}

func groupMsgRejectHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupMsgRejectData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupMsgReject != nil {
		return DefaultHandlers.GroupMsgReject(payload, data)
	}
	return nil
}

func groupMsgReceiveHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGroupMsgReceiveData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.GroupMsgReceive != nil {
		return DefaultHandlers.GroupMsgReceive(payload, data)
	}
	return nil
}

func friendAddHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSFriendAddData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.FriendAdd != nil {
		return DefaultHandlers.FriendAdd(payload, data)
	}
	return nil
}

func friendDelHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSFriendDelData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.FriendDel != nil {
		return DefaultHandlers.FriendDel(payload, data)
	}
	return nil
}

func c2cMsgRejectHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSFriendMsgRejectData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.C2CMsgReject != nil {
		return DefaultHandlers.C2CMsgReject(payload, data)
	}
	return nil
}

func c2cMsgReceiveHandle(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSFriendMsgReveiceData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.C2CMsgReceive != nil {
		return DefaultHandlers.C2CMsgReceive(payload, data)
	}
	return nil
}