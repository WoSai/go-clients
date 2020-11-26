package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jacexh/requests"
)

type (
	Client struct {
		url    string
		client *requests.Session
		option Option
		ak     *AccessTokenKeeper
	}

	AccessTokenKeeper struct {
		mu sync.RWMutex
		ak string
	}
)

func (ak *AccessTokenKeeper) Set(token string) {
	ak.mu.Lock()
	defer ak.mu.Unlock()
	ak.ak = token
}

func (ak *AccessTokenKeeper) Get() string {
	ak.mu.RLock()
	defer ak.mu.RUnlock()
	return ak.ak
}

func NewClient(opt Option) *Client {
	return &Client{
		url: "https://oapi.dingtalk.com",
		client: requests.NewSession(requests.Option{
			Name:    "github.com/wosai/go-clients/dingtalk",
			Timeout: 30 * time.Second,
		}),
		option: opt,
		ak:     &AccessTokenKeeper{},
	}
}

// GetUserInfoByCode 根据sns临时授权码获取用户信息 https://ding-doc.dingtalk.com/document#/org-dev-guide/obtain-the-user-information-based-on-the-sns-temporary-authorization#topic-1995619
func (client *Client) GetUserInfoByCode(ctx context.Context, code *RequestGetUserInfoByCode) (*UserInfo, *http.Response, error) {
	ts := time.Now().UnixNano() / 1e6
	h := hmac.New(sha256.New, []byte(client.option.AppSecret))
	h.Write([]byte(strconv.FormatInt(ts, 10)))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	ret := new(ResponseGetUserInfoByCode)
	res, _, err := client.client.PostWithContext(
		ctx,
		client.url+"/sns/getuserinfo_bycode",
		requests.Params{
			Query: requests.Any{"accessKey": client.option.AppID, "timestamp": strconv.FormatInt(ts, 10), "signature": signature},
			Json:  code,
		},
		requests.UnmarshalJSONResponse(ret))
	if err != nil {
		return nil, res, err
	}
	return ret.UserInfo, res, ret.DingtalkErr.Err()
}

// GetAccessToken 获取access_token https://ding-doc.dingtalk.com/document#/org-dev-guide/obtain-access_token
func (client *Client) GetAccessToken(ctx context.Context) (string, *http.Response, error) {
	ret := new(ResponseGetAccessToken)
	res, _, err := client.client.GetWithContext(
		ctx,
		client.url+"/gettoken",
		requests.Params{Query: requests.Any{"appkey": client.option.AppID, "appsecret": client.option.AppSecret}},
		requests.UnmarshalJSONResponse(ret),
	)
	if err != nil {
		return "", res, err
	}
	if ret.DingtalkErr.Err() == nil {
		client.ak.Set(ret.AccessToken)
	}
	return ret.AccessToken, res, ret.DingtalkErr.Err()
}

func (client *Client) RetryOnAccessTokenExpired() {}
