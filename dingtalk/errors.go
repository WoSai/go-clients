package dingtalk

import "fmt"

// DingtalkErr 钉钉错误信息
type DingtalkErr struct {
	ErrorMessage string `json:"errmsg,omitempty"`
	ErrorCode    int    `json:"errcode,omitempty"`
	SubCode      string `json:"sub_code,omitempty"`
	SubMessage   string `json:"sub_msg,omitempty"`
}

// GotErr 返回response中的错误
func (de *DingtalkErr) GotErr() error {
	if de == nil || de.ErrorMessage == "ok" || de.ErrorCode == 0 {
		return nil
	}
	return de
}

// Error error的实现
func (de *DingtalkErr) Error() string {
	if de.ErrorCode == AuthenticationAbnormal {
		return fmt.Sprintf("[%s]:%s", de.SubCode, de.SubMessage)
	}
	return fmt.Sprintf("[%d]: %s", de.ErrorCode, de.ErrorMessage)
}

// IsAccessTokenExpired access_token是否过期
func (de *DingtalkErr) IsAccessTokenExpired() bool {
	if de.GotErr() != nil &&
		((de.ErrorCode == AuthenticationAbnormal &&
			(de.SubCode == InvalidAccessToken || de.SubCode == EmptyAccessToken || de.SubCode == IllegalAccessToken)) ||
			de.ErrorCode == 40014) {
		return true
	}
	return false
}

const (
	// https://ding-doc.dingtalk.com/document#/org-dev-guide/server-api-error-codes
	// AuthenticationAbnormal 鉴权异常
	AuthenticationAbnormal = 88
	// InvalidAccessToken 获取access_token时Secret错误，或者access_token无效
	InvalidAccessToken = "40001"
	EmptyAccessToken   = "40000"
	IllegalAccessToken = "40014"
)
