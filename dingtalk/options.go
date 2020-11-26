package dingtalk

type (
	// Option 应用凭证
	Option struct {
		AgentID   string
		AppKey    string
		AppSecret string
	}
)

func (ao Option) IsEmpty() bool {
	if ao.AppKey == "" || ao.AppSecret == "" {
		return true
	}
	return false
}
