package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jacexh/requests"
	"golang.org/x/sync/singleflight"
)

type (
	Client struct {
		mu     sync.RWMutex
		url    string
		client *requests.Session
		opt    Option
		flight singleflight.Group
		ak     string
	}
)

func NewClient(opt Option) *Client {
	return &Client{
		url: "https://oapi.dingtalk.com",
		client: requests.NewSession(requests.Option{
			Name:    "github.com/wosai/go-clients/dingtalk",
			Timeout: 30 * time.Second,
		}),
		opt: opt,
	}
}

func (ding *Client) WithAppOption(opt Option) *Client {
	ding.opt = opt
	return ding
}

func (ding *Client) AccessToken() string {
	ding.mu.RLock()
	defer ding.mu.RUnlock()
	return ding.ak
}

func (ding *Client) SetAccessToken(token string) {
	ding.mu.Lock()
	defer ding.mu.Unlock()
	ding.ak = token
}

// GetUserInfoByCode 根据sns临时授权码获取用户信息 https://ding-doc.dingtalk.com/document#/org-dev-guide/obtain-the-user-information-based-on-the-sns-temporary-authorization#topic-1995619
func (ding *Client) GetUserInfoByCode(ctx context.Context, code *RequestGetUserInfoByCode) (*UserInfo, *http.Response, error) {
	var err error
	var res *http.Response
	ret := new(ResponseGetUserInfoByCode)

	ts := time.Now().UnixNano() / 1e6
	h := hmac.New(sha256.New, []byte(ding.opt.LoginAppSecret))
	h.Write([]byte(strconv.FormatInt(ts, 10)))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	res, _, err = ding.client.PostWithContext(
		ctx,
		ding.url+"/sns/getuserinfo_bycode",
		requests.Params{
			Query: requests.Any{"accessKey": ding.opt.LoginAppID, "timestamp": strconv.FormatInt(ts, 10), "signature": signature},
			Json:  code,
		},
		UnmarshalAndParseError(ret))

	if err == nil {
		return ret.UserInfo, res, nil
	}
	return nil, res, err
}

// GetAccessToken 获取access_token https://ding-doc.dingtalk.com/document#/org-dev-guide/obtain-access_token
func (ding *Client) GetAccessToken(ctx context.Context) (string, *http.Response, error) {
	if ding.opt.IsEmpty() {
		return "", nil, errors.New("no app provided")
	}
	ret := new(ResponseGetAccessToken)
	res, _, err := ding.client.GetWithContext(
		ctx,
		ding.url+"/gettoken",
		requests.Params{Query: requests.Any{"appkey": ding.opt.AppKey, "appsecret": ding.opt.AppSecret}},
		UnmarshalAndParseError(ret),
	)
	if err != nil {
		return "", res, err
	}
	ding.SetAccessToken(ret.AccessToken)
	return ret.AccessToken, res, nil
}

func (ding *Client) GetUserByUnionID(ctx context.Context, req *RequestGetByUnionID) (*UserGetByUnionId, *http.Response, error) {
	ret := new(ResponseGetByUnionID)
	var res *http.Response
	var err error

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/user/getbyunionid",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: req},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	if err == nil {
		return ret.Result, res, nil
	}
	return nil, res, err
}

func (ding *Client) RetryOnAccessTokenExpired(ctx context.Context, retry int, fn func() error) (err error) {
	for i := 0; i < retry+1; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		var de *DingtalkErr
		if errors.As(err, &de) && err.(*DingtalkErr).IsAccessTokenExpired() {
			_, akErr, _ := ding.flight.Do("access_token", func() (interface{}, error) {
				ak, _, reqErr := ding.GetAccessToken(ctx)
				ding.flight.Forget("access_token")
				return ak, reqErr
			})
			if akErr != nil {
				return fmt.Errorf("%w | %s", err, akErr.Error())
			} else {
				continue
			}
		} else {
			return err
		}
	}
	return err
}

// GetOrganizationUserCount 获取组织人员数
func (ding *Client) GetOrganizationUserCount(ctx context.Context, onlyActive int) (int, *http.Response, error) {
	var res *http.Response
	var err error
	ret := new(ResponseOrganizationUserCount)

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.GetWithContext(
			ctx,
			ding.url+"/user/get_org_user_count",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken(), "onlyActive": strconv.Itoa(onlyActive)}},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	if err == nil {
		return ret.Count, res, nil
	}
	return 0, res, err
}

// GetUser 获取用户 https://ding-doc.dingtalk.com/document#/org-dev-guide/queries-user-details#topic-1960047
func (ding *Client) GetUser(ctx context.Context, get *RequestUserGet) (*UserGetRequest, *http.Response, error) {
	var res *http.Response
	var err error
	ret := new(ResponseUserGet)

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/v1/user/get",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: get},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	if err != nil {
		return nil, res, err
	}
	return ret.Result, res, err
}
