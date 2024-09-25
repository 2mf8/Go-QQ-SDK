package proxy

import "net/http"

var Proxy = http.ProxyFromEnvironment
