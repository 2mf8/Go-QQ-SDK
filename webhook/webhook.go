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

	"github.com/2mf8/Go-QQ-SDK/dto"
	"github.com/2mf8/Go-QQ-SDK/openapi"
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
	FirstStart bool = true
)

var Bots = make(map[string]*Bot)

type Bot struct {
	QQ        uint64
	AppId     string
	Token     string
	AppSecret string
	Openapi   openapi.OpenAPI

	mux           sync.RWMutex
	WaitingFrames map[string]*promise.Promise

	Payload *dto.WSPayload
}

type BotHeaderInfo struct {
	ContentLength       []string `json:"Content-Length,omitempty"`
	ContentType         []string `json:"Content-Type,omitempty"`
	UserAgent           []string `json:"User-Agent,omitempty"`
	XBotAppid           []string `json:"X-Bot-Appid,omitempty"`
	XSignatureEd25519   []string `json:"X-Signature-Ed25519,omitempty"`
	XSignatureMethod    []string `json:"X-Signature-Method,omitempty"`
	XSignatureTimestamp []string `json:"X-Signature-Timestamp,omitempty"`
	XTpsTraceId         []string `json:"X-Tps-Trace-Id,omitempty"`
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
	appid := c.Param("appid")
	//if appid == "" {} else {}
	fmt.Println(appid)
	header := &BotHeaderInfo{}
	h, _ := json.Marshal(c.Request.Header)
	json.Unmarshal(h, header)
	fmt.Println("Header信息：", string(h))
	httpBody, err := io.ReadAll(c.Request.Body)
	fmt.Println("Body信息：", string(httpBody))
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
	if FirstStart {
		NewBot(header, payload, b, header.XBotAppid[0])
		FirstStart = false
	}
	NewBot(header, payload, b, header.XBotAppid[0])
	if err = json.Unmarshal(b, validationPayload); err != nil {
		log.Println("parse http payload failed:", err)
		return
	}
	seed := AllSetting.Apps[appid].AppSecret
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
	router.POST("/qqbot/:appid", handleValidation)

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
		if AllSetting.CertFile == "" || AllSetting.CertKey == "" {
			if err := http.Serve(ln, engine); err != nil {
				FatalError(fmt.Errorf("failed to serve http, err: %+v", err))
			}
		} else {
			if err := http.ServeTLS(ln, engine, AllSetting.CertFile, AllSetting.CertKey); err != nil {
				FatalError(fmt.Errorf("failed to serve http, err: %+v", err))
			}
		}
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

func NewBot(h *BotHeaderInfo, p *dto.WSPayload, m []byte, appId string) *Bot {
	as := ReadSetting()
	ibot, ok := Bots[appId]
	if ok {
		ibot.ParseWHData(h, p, m)
	}
	bot := &Bot{
		AppId:     appId,
		Token:     as.Apps[appId].Token,
		AppSecret: as.Apps[appId].AppSecret,
		Payload:   p,
	}
	Bots[bot.AppId] = bot
	return bot
}

func (bot *Bot) AddOpenapi(iOpenapi openapi.OpenAPI) *Bot {
	bot.Openapi = iOpenapi
	return bot
}

func (bot *Bot) ParseWHData(h *BotHeaderInfo, p *dto.WSPayload, message []byte) {
	if p.Type == dto.EventGroupATMessageCreate {
		gm := &dto.WSGroupATMessageData{}
		err := json.Unmarshal(message, gm)
		if err == nil {
			GroupAtMessageEventHandler(h, p, gm)
		}
	}
	if p.Type == dto.EventGroupAddRobbot {
		gar := &dto.WSGroupAddRobotData{}
		err := json.Unmarshal(message, gar)
		if err == nil {
			GroupAddRobotEventHandler(h, p, gar)
		}
	}
	if p.Type == dto.EventGroupDelRobbot {
		gdr := &dto.WSGroupDelRobotData{}
		err := json.Unmarshal(message, gdr)
		if err == nil {
			GroupDelRobotEventHandler(h, p, gdr)
		}
	}
	if p.Type == dto.EventGroupMsgReceive {
		gmr := &dto.WSGroupMsgReceiveData{}
		err := json.Unmarshal(message, gmr)
		if err == nil {
			GroupMsgReceiveEventHandler(h, p, gmr)
		}
	}
	if p.Type == dto.EventGroupMsgReject {
		gmr := &dto.WSGroupMsgRejectData{}
		err := json.Unmarshal(message, gmr)
		if err == nil {
			GroupMsgRejectEventHandler(h, p, gmr)
		}
	}
	if p.Type == dto.EventC2CMessageCreate {
		cmc := &dto.WSC2CMessageData{}
		err := json.Unmarshal(message, cmc)
		if err == nil {
			C2CMessageEventHandler(h, p, cmc)
		}
	}
	if p.Type == dto.EventC2CMsgReceive {
		fmr := &dto.WSFriendMsgReveiceData{}
		err := json.Unmarshal(message, fmr)
		if err == nil {
			C2CMsgReceiveHandler(h, p, fmr)
		}
	}
	if p.Type == dto.EventC2CMsgReject {
		fmr := &dto.WSFriendMsgRejectData{}
		err := json.Unmarshal(message, fmr)
		if err == nil {
			C2CMsgRejectHandler(h, p, fmr)
		}
	}
	if p.Type == dto.EventFriendAdd {
		fad := &dto.WSFriendAddData{}
		err := json.Unmarshal(message, fad)
		if err == nil {
			FriendAddEventHandler(h, p, fad)
		}
	}
	if p.Type == dto.EventFriendDel {
		fad := &dto.WSFriendDelData{}
		err := json.Unmarshal(message, fad)
		if err == nil {
			FriendDelEventHandler(h, p, fad)
		}
	}
	if p.Type == dto.EventAtMessageCreate {
		am := &dto.WSATMessageData{}
		err := json.Unmarshal(message, am)
		if err == nil {
			ATMessageEventHandler(h, p, am)
		}
	}
	if p.Type == dto.EventMessageCreate {
		m := &dto.WSMessageData{}
		err := json.Unmarshal(message, m)
		if err == nil {
			MessageEventHandler(h, p, m)
		}
	}
	if p.Type == dto.EventInteractionCreate {
		i := &dto.WSInteractionData{}
		err := json.Unmarshal(message, i)
		if err == nil {
			i.ID = p.ID
			InteractionEventHandler(h, p, i)
		}
	}
	if p.Type == dto.EventDirectMessageCreate {
		i := &dto.WSDirectMessageData{}
		err := json.Unmarshal(message, i)
		if err == nil {
			DirectMessageEventHandler(h, p, i)
		}
	}
	if p.Type == dto.EventMessageReactionAdd || p.Type == dto.EventMessageReactionRemove {
		mr := &dto.WSMessageReactionData{}
		err := json.Unmarshal(message, mr)
		if err == nil {
			MessageReactionEventHandler(h, p, mr)
		}
	}
	if p.Type == dto.EventMessageAuditPass || p.Type == dto.EventMessageAuditReject {
		mr := &dto.WSMessageAuditData{}
		err := json.Unmarshal(message, mr)
		if err == nil {
			MessageAuditEventHandler(h, p, mr)
		}
	}
	if p.Type == dto.EventForumThreadCreate || p.Type == dto.EventForumPostCreate || p.Type == dto.EventForumReplyCreate || p.Type == dto.EventForumThreadUpdate || p.Type == dto.EventForumPostDelete || p.Type == dto.EventForumThreadDelete || p.Type == dto.EventForumReplyDelete {
		ft := &dto.WSForumAuditData{}
		err := json.Unmarshal(message, ft)
		if err == nil {
			ForumAuditEventHandler(h, p, ft)
		}
	}
	if p.Type == dto.EventGuildCreate || p.Type == dto.EventGuildUpdate || p.Type == dto.EventGuildDelete {
		g := &dto.WSGuildData{}
		err := json.Unmarshal(message, g)
		if err == nil {
			GuildEventHandler(h, p, g)
		}
	}
	if p.Type == dto.EventChannelCreate || p.Type == dto.EventChannelUpdate || p.Type == dto.EventChannelDelete {
		c := &dto.WSChannelData{}
		err := json.Unmarshal(message, c)
		if err == nil {
			ChannelEventHandler(h, p, c)
		}
	}
	if p.Type == dto.EventGuildMemberAdd || p.Type == dto.EventGuildMemberUpdate || p.Type == dto.EventGuildMemberRemove {
		gm := &dto.WSGuildMemberData{}
		err := json.Unmarshal(message, gm)
		if err == nil {
			GuildMemberEventHandler(h, p, gm)
		}
	}
}
