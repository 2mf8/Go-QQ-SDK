package bot

import (
	"net/http"
	"net/url"

	"github.com/2mf8/Go-QQ-Client/internal/proxy"
	"github.com/2mf8/Go-QQ-Client/log"
	"github.com/2mf8/Go-QQ-Client/openapi"
	"github.com/2mf8/Go-QQ-Client/websocket"
)

// SetLogger 设置 logger，需要实现 sdk 的 log.Logger 接口
func SetLogger(logger log.Logger) {
	log.DefaultLogger = logger
}

// SetSessionManager 注册自己实现的 session manager
func SetSessionManager(m SessionManager) {
	defaultSessionManager = m
}

// SetWebsocketClient 替换 websocket 实现
func SetWebsocketClient(c websocket.WebSocket) {
	websocket.Register(c)
}

// SetOpenAPIClient 注册 openapi 的不同实现，需要设置版本
func SetOpenAPIClient(v openapi.APIVersion, c openapi.OpenAPI) {
	openapi.Register(v, c)
}

// SetProxy 设置代理
//
// 若调用了SetWebsocketClient，会导致websocket代理无效，请自行设置websocket代理
//
// 若无设置，则使用系统代理
func SetProxy(p func(*http.Request) (*url.URL, error)) {
	proxy.Proxy = p
}
