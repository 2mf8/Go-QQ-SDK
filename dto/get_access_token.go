package dto

type GetAccessTokenReq struct {
	AppId        uint64 `json"appId,omitempty"`
	ClientSecret string `json"clientSecret,omitempty"`
}

type GetAccessTokenResp struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"`
}
