package dockerhub

type Option struct {
	URL       string // 为空时表示使用官方仓库
	AuthToken string // 有则优先使用
	Username  string // 用户名
	Password  string // 密码
}
