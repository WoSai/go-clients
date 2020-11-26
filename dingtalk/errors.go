package dingtalk

import "fmt"

// DingtalkErr 钉钉错误信息
type DingtalkErr struct {
	ErrorMessage string `json:"errmsg,omitempty"`
	ErrorCode    int    `json:"errcode,omitempty"`
}

func (de *DingtalkErr) Err() error {
	if de == nil || de.ErrorMessage == "ok" || de.ErrorCode == 0 {
		return nil
	}
	return fmt.Errorf("[%d]: %s", de.ErrorCode, de.ErrorMessage)
}
