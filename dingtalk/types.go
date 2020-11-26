package dingtalk

import (
	"encoding/json"
	"time"
)

type (
	UserInfo struct {
		Nick                 string `json:"nick,omitempty"`                     // 用户在钉钉上面的昵称
		UnionID              string `json:"unionid,omitempty"`                  // 用户在当前开放应用所属企业的唯一标识
		OpenID               string `json:"openid,omitempty"`                   // 用户在当前开放应用内的唯一标识
		MainOrgAuthHighLevel bool   `json:"main_org_auth_high_level,omitempty"` // 用户主企业是否达到高级认证级别
	}

	BasicResponse struct {
		RequestID    string `json:"request_id,omitempty"`
		*DingtalkErr `json:",inline"`
	}

	// ResponseGetUserInfoByCode https://oapi.dingtalk.com/sns/getuserinfo_bycode
	ResponseGetUserInfoByCode struct {
		UserInfo     *UserInfo `json:"user_info"`
		*DingtalkErr `json:",inline"`
	}

	// RequestGetUserInfoByCode https://oapi.dingtalk.com/sns/getuserinfo_bycode
	RequestGetUserInfoByCode struct {
		TempAuthCode string `json:"tmp_auth_code"`
	}

	// RequestGetByUnionID https://oapi.dingtalk.com/topapi/user/getbyunionid
	RequestGetByUnionID struct {
		UnionID string `json:"unionid"`
	}

	UserGetByUnionId struct {
		ContactType int    `json:"contact_type,omitempty"`
		UserID      string `json:"userid,omitempty"`
	}

	// ResponseGetByUnionID https://oapi.dingtalk.com/topapi/user/getbyunionid
	ResponseGetByUnionID struct {
		*BasicResponse `json:",inline"`
		Result         *UserGetByUnionId `json:"result,omitempty"`
	}

	// DepartmentOrder 部门排序
	DepartmentOrder struct {
		DepartmentID int `json:"dept_id"`
		Order        int `json:"order"`
	}

	DepartmentLeader struct {
		DepartmentID int  `json:"dept_id"`
		Leader       bool `json:"leader"`
	}

	UserRole struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Group string `json:"group_name"`
	}

	UnixTimestamp struct {
		ts int64
	}

	UserGetRequest struct {
		UserID             string              `json:"userid"`
		UnionID            string              `json:"unionid"`
		Name               string              `json:"name"`
		Avatar             string              `json:"avatar"`
		StateCode          string              `json:"state_code"` // 国际电话区号
		Mobile             string              `json:"mobile"`
		HideMobile         bool                `json:"hide_mobile"` // 是否隐藏手机号码
		Telephone          string              `json:"telephone"`   // 分机号
		JobNumber          string              `json:"job_number"`  // 工号
		Title              string              `json:"title"`       // 职位
		Email              string              `json:"email"`       // 邮箱
		WorkPlace          string              `json:"work_place"`  // 办公地点
		Remark             string              `json:"remark"`
		DepartmentIDs      []int               `json:"dept_id_list"`    // 所属部门ID列表
		DepartmentOrders   []*DepartmentOrder  `json:"dept_order_list"` // 员工在对应部门的排序
		Extension          string              `json:"extension"`       // 扩展属性
		HiredDate          *UnixTimestamp      `json:"hired_date"`
		Active             bool                `json:"active"`
		RealAuthed         bool                `json:"real_authed"` // 是否完成了实名认证
		Admin              bool                `json:"admin"`       // 是否未企业管理员
		Boss               bool                `json:"boss"`        // 是否为企业老板
		LeaderInDepartment []*DepartmentLeader `json:"leader_in_dept"`
		Roles              []*UserRole         `json:"role_list"`
	}

	// RequestUserGet https://oapi.dingtalk.com/topapi/v2/user/get
	RequestUserGet struct {
		UserID   string `json:"userid"`
		Language string `json:"language,omitempty"`
	}

	// ResponseUserGet https://oapi.dingtalk.com/topapi/v2/user/get
	ResponseUserGet struct {
		BasicResponse `json:",inline"`
		Result        *UserGetRequest `json:"result"`
	}

	ResponseGetAccessToken struct {
		*DingtalkErr `json:",inline"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	ResponseOrganizationUserCount struct {
		*DingtalkErr `json:",inline"`
		Count        int `json:"count"`
	}
)

func (ts *UnixTimestamp) UnmarshalJSON(data []byte) error {
	var t int64
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	ts.ts = t
	return nil
}

func (ts *UnixTimestamp) Time() time.Time {
	return time.Unix(0, ts.ts*1e6)
}

func (ts *UnixTimestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(ts.ts)
}
