package webhook

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/2mf8/Go-QQ-Client/openapi"
	"github.com/fanliao/go-promise"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"google.golang.org/protobuf/proto"
)

var (
	SelectPort = map[string]string{
		"80":   ":80",
		"8080": ":8080",
		"443":  ":443",
		"8443": ":8443",
	}
)

// WSGuildData 频道 payload
var WHGuildData = &dto.Guild{}

// WSGuildMemberData 频道成员 payload
var WHGuildMemberData = &dto.Member{}

// WSChannelData 子频道 payload
var WHChannelData = &dto.Channel{}

// WSMessageData 消息 payload
var WHMessageData = &dto.Message{}

// WSATMessageData only at 机器人的消息 payload
var WHATMessageData = &dto.Message{}

// WSDirectMessageData 私信消息 payload
var WHDirectMessageData = &dto.Message{}

var WHC2CMessageData = &dto.C2CMessage{}

var WHGroupATMessageData = &dto.GroupMessage{}

var WHGroupMessageData = &dto.GroupMessage{}

// WSMessageDeleteData 消息 payload
var WHMessageDeleteData = &dto.MessageDelete{}

// WSPublicMessageDeleteData 公域机器人的消息删除 payload
var WHPublicMessageDeleteData = &dto.MessageDelete{}

// WSDirectMessageDeleteData 私信消息 payload
var WHDirectMessageDeleteData = &dto.MessageDelete{}

// WSAudioData 音频机器人的音频流事件
var WHAudioData = &dto.AudioAction{}

// WSMessageReactionData 表情表态事件
var WHMessageReactionData = &dto.MessageReaction{}

// WSMessageAuditData 消息审核事件
var WHMessageAuditData = &dto.MessageAudit{}

// WSThreadData 主题事件
var WHThreadData = &dto.Thread{}

// WSPostData 帖子事件
var WHPostData = &dto.Post{}

// WSReplyData 帖子回复事件
var WHReplyData = &dto.Reply{}

// WSForumAuditData 帖子审核事件
var WHForumAuditData = &dto.ForumAuditResult{}

// WSInteractionData 互动事件
var WHInteractionData = &dto.Interaction{}

var WHGroupAddRobotData = &dto.GroupAddRobotEvent{}

var WHGroupDelRobotData = &dto.GroupDelRobotEvent{}

var WHGroupMsgRejectData = &dto.GroupMsgRejectEvent{}

var WHGroupMsgReceiveData = &dto.GroupMsgReceiveEvent{}

var WHFriendAddData = &dto.FriendAddEvent{}

var WHFriendDelData = &dto.FriendDelEvent{}

var WHFriendMsgRejectData = &dto.FriendMsgRejectEvent{}

var WHFriendMsgReveiceData = &dto.FriendMsgReceiveEvent{}

var Bots = make(map[int64]*Bot)

type Bot struct {
	QQ        uint64
	AppId     uint64
	Token     string
	AppSecret string
	Openapi   openapi.OpenAPI

	mux           sync.RWMutex
	WaitingFrames map[string]*promise.Promise

	Payload *dto.WSPayload
}
type ValidationRequest struct {
	PlainToken string `json:"plain_token,omitempty"`
	EventTs    string `json:"event_ts,omitempty"`
}

type ValidationResponse struct {
	PlainToken string `json:"plain_token,omitempty"`
	Signature  string `json:"signature,omitempty"`
}

func handleValidation(c *gin.Context) {
	httpBody, err := io.ReadAll(c.Request.Body)
	fmt.Println(string(httpBody))
	if err != nil {
		log.Println("read http body err", err)
		return
	}
	payload := &dto.WSPayload{}
	if err = json.Unmarshal(httpBody, payload); err != nil {
		log.Println("parse http payload err", err)
		return
	}
	validationPayload := &ValidationRequest{}
	b, _ := json.Marshal(payload.Data)
	NewBot(payload, b, AllSetting.AppId)
	if err = json.Unmarshal(b, validationPayload); err != nil {
		log.Println("parse http payload failed:", err)
		return
	}
	seed := AllSetting.AppSecret
	for len(seed) < ed25519.SeedSize {
		seed = strings.Repeat(seed, 2)
	}
	seed = seed[:ed25519.SeedSize]
	reader := strings.NewReader(seed)
	// GenerateKey 方法会返回公钥、私钥，这里只需要私钥进行签名生成不需要返回公钥
	_, privateKey, err := ed25519.GenerateKey(reader)
	if err != nil {
		log.Println("ed25519 generate key failed:", err)
		return
	}
	var msg bytes.Buffer
	msg.WriteString(validationPayload.EventTs)
	msg.WriteString(validationPayload.PlainToken)
	signature := hex.EncodeToString(ed25519.Sign(privateKey, msg.Bytes()))
	if err != nil {
		log.Println("generate signature failed:", err)
		return
	}
	rspBytes, err := json.Marshal(
		&ValidationResponse{
			PlainToken: validationPayload.PlainToken,
			Signature:  signature,
		})
	if err != nil {
		log.Println("handle validation failed:", err)
		return
	}
	c.Data(http.StatusOK, c.ContentType(), rspBytes)
}

func InitGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(CORSMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "it works")
	})
	router.POST("/qqbot", handleValidation)

	iport := strconv.FormatInt(int64(AllSetting.Port), 10)
	realPort, err := RunGin(router, ":"+iport)
	if err != nil {
		for i, v := range SelectPort {
			if i == iport {
				continue
			} else {
				iport = i
				realPort, err := RunGin(router, v)
				if err != nil {
					log.Warn(fmt.Errorf("failed to run gin, err: %+v", err))
					continue
				}
				iport = realPort
				log.Infof("端口号 %s", realPort)
				break
			}
		}
	} else {
		iport = realPort
		log.Infof("端口号 %s", realPort)
	}
}

func RunGin(engine *gin.Engine, port string) (string, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return "", err
	}
	_, randPort, _ := net.SplitHostPort(ln.Addr().String())
	go func() {
		/* if AllSetting.CertFile == "" || AllSetting.CertKey == "" { */
		if err := http.Serve(ln, engine); err != nil {
			FatalError(fmt.Errorf("failed to serve http, err: %+v", err))
		}
		/* }else{
			if err := http.ServeTLS(ln, engine, AllSetting.CertFile, AllSetting.CertKey); err != nil {
				FatalError(fmt.Errorf("failed to serve http, err: %+v", err))
			}
		} */
	}()
	return randPort, nil
}

func InitLog() {
	// 输出到命令行
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	// 输出到文件
	rotateLogs, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join("logs", "latest.log")), // 最新日志软链接
		rotatelogs.WithRotationTime(time.Hour*24),                // 每天一个新文件
		rotatelogs.WithMaxAge(time.Hour*24*3),                    // 日志保留3天
	)
	if err != nil {
		FatalError(err)
		return
	}
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  rotateLogs,
			log.WarnLevel:  rotateLogs,
			log.ErrorLevel: rotateLogs,
			log.FatalLevel: rotateLogs,
			log.PanicLevel: rotateLogs,
		},
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%time%] [%lvl%]: %msg% \r\n",
		},
	))
}

func FatalError(err error) {
	log.Errorf(err.Error())
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	sBuf := string(buf)
	log.Errorf(sBuf)
	time.Sleep(5 * time.Second)
	panic(err)
}

func Return(c *gin.Context, resp proto.Message) {
	var (
		data []byte
		err  error
	)
	switch c.ContentType() {
	case binding.MIMEPROTOBUF:
		data, err = proto.Marshal(resp)
	case binding.MIMEJSON:
		data, err = json.Marshal(resp)
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "marshal resp error")
		return
	}
	c.Data(http.StatusOK, c.ContentType(), data)
}

func NewBot(p *dto.WSPayload, m []byte, appId uint64) *Bot {
	as := ReadSetting()
	ibot, ok := Bots[int64(appId)]
	if ok {
		ibot.ParseWHData(p, m)
	}
	bot := &Bot{
		AppId:     as.AppId,
		Token:     as.Token,
		AppSecret: as.AppSecret,
		Payload:   p,
	}
	Bots[int64(bot.AppId)] = bot
	return bot
}

func (bot *Bot) AddOpenapi(iOpenapi openapi.OpenAPI) *Bot {
	bot.Openapi = iOpenapi
	return bot
}

func (bot *Bot) ParseWHData(p *dto.WSPayload, message []byte) {
	if p.Type == dto.EventGroupATMessageCreate {
		gm := &dto.WSGroupATMessageData{}
		err := json.Unmarshal(message, gm)
		if err == nil {
			GroupAtMessageEventHandler(p, gm)
		}
	}
	if p.Type == dto.EventGroupAddRobbot {
		gar := &dto.WSGroupAddRobotData{}
		err := json.Unmarshal(message, gar)
		if err == nil {
			GroupAddRobotEventHandle(p, gar)
		}
	}
}
