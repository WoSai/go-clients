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

// 获取用户详情V2 https://ding-doc.dingtalk.com/document#/org-dev-guide/queries-user-details
func (ding *Client) GetUserInfoV2(ctx context.Context, req *RequestUserGet) (*ResponseGetUserInfo, *http.Response, error) {
	var res *http.Response
	var err error
	ret := new(ResponseGetUserInfo)

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/v2/user/get",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: req},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret, res, err
}

// 根据手机号获取userid v2  https://ding-doc.dingtalk.com/document#/org-dev-guide/query-users-by-phone-number
func (ding *Client) GetUserInfoByMobileV2(ctx context.Context, mobile string) (string, *http.Response, error) {
	ret := new(ResponseUserByMobile)
	var res *http.Response
	var err error

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/v2/user/getbymobile",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: requests.Any{"mobile": mobile}},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	if err == nil {
		return ret.Result.Userid, res, err
	}
	return "", res, err
}

// 获取指定用户的所有父部门列表  https://ding-doc.dingtalk.com/document#/org-dev-guide/obtains-the-list-of-all-parent-departments-of-a-user
func (ding *Client) ListParentDeptByUserV2(ctx context.Context, userid string) (*DeptListParent, *http.Response, error) {
	ret := new(ResponseGetUserIdByUnionid)
	var res *http.Response
	var err error

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/v2/department/listparentbyuser",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: requests.Any{"userid": userid}},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret.Result, res, err
}

// 获取部门详情 https://ding-doc.dingtalk.com/document#/org-dev-guide/queries-department-details-v1
func (ding *Client) GetDepartment(ctx context.Context, req *RequestDepartmentInfo) (*DepartmentInfo, *http.Response, error) {
	ret := new(ResponseDeptInfo)
	var res *http.Response
	var err error

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		query := requests.Any{"access_token": ding.AccessToken(), "id": req.DeptId}
		if req.Language != "" {
			query["lang"] = string(req.Language)
		}
		res, _, err = ding.client.GetWithContext(
			ctx,
			ding.url+"/department/get",
			requests.Params{Query: query},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret.DepartmentInfo, res, err
}

// 发起审批实例 https://developers.dingtalk.com/document/app/initiate-approval
func (ding *Client) CreateProcessInstance(ctx context.Context, req *RequestCreateProcessInstance) (string, *http.Response, error) {
	ret := new(ResponseCreateProcessInstance)
	var err error
	var res *http.Response

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/processinstance/create",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: req},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret.ProcessInstanceID, res, err
}

// 获取审批实例详情 https://developers.dingtalk.com/document/app/obtains-the-details-of-a-single-approval-instance
func (ding *Client) GetProcessInstance(ctx context.Context, processInstanceID string) (*ProcessInstance, *http.Response, error) {
	ret := new(ResponseGetProcessInstance)
	var err error
	var res *http.Response

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/processinstance/get",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: requests.Any{"process_instance_id": processInstanceID}},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret.ProcessInstance, res, err
}

// 获取部门详情,可以得到部门主管
func (ding *Client) GetDepartmentV2(ctx context.Context, deptID int) (*DepartmentInfoV2, *http.Response, error) {
	ret := new(ResponseDeptInfoV2)
	var res *http.Response
	var err error

	err = ding.RetryOnAccessTokenExpired(ctx, 1, func() error {
		res, _, err = ding.client.PostWithContext(
			ctx,
			ding.url+"/topapi/v2/department/get",
			requests.Params{Query: requests.Any{"access_token": ding.AccessToken()}, Json: &RequestDepartmentInfoV2{DeptID: deptID}},
			UnmarshalAndParseError(ret),
		)
		return err
	})
	return ret.Result, res, err
}
