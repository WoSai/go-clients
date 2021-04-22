package dingtalk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

var DingClient *Client
var ctx = context.Background()
var testUser *ResponseGetUserInfo
var DeptID int

var (
	AgentID   = os.Getenv("AgentID")
	AppKey    = os.Getenv("AppKey")
	AppSecret = os.Getenv("AppSecret")
	UserID    = os.Getenv("UserID")
	_deptID    = os.Getenv("DeptID")
)

func init() {
	var err error
	DeptID, err = strconv.Atoi(_deptID)
	if err != nil {
		panic(fmt.Errorf("%w DeptID: %s not int", err, _deptID))
	}
	DingClient = NewClient(Option{AgentID: AgentID, AppKey: AppKey, AppSecret: AppSecret})
	_, _, err = DingClient.GetAccessToken(ctx)
	if err != nil {
		panic("get access token fail:" + err.Error())
	}
	testUser, _, err = DingClient.GetUserInfoV2(ctx, &RequestUserGet{UserID: UserID})
	if err != nil {
		panic("userid not exit" + err.Error())
	}
}

func TestClient_WithAppOption(t *testing.T) {
	opt := Option{AgentID: "AgentID", AppKey: "AppKey", AppSecret: "AppSecret"}
	DingClient.WithAppOption(opt)
	assert.Equal(t, DingClient.opt, opt)
}

func TestClient_GetUserInfoByCode(t *testing.T) {
	_, _, err := DingClient.GetUserInfoByCode(ctx, &RequestGetUserInfoByCode{TempAuthCode: "testcode"})
	assert.NotNil(t, err)
}

func TestClient_GetUserByUnionID(t *testing.T) {
	user, _, err := DingClient.GetUserByUnionID(ctx, &RequestGetByUnionID{UnionID: testUser.Result.UnionID})
	assert.Nil(t, err)
	assert.Equal(t, user.UserID, testUser.Result.UserID)
}

func TestClient_GetDepartment(t *testing.T) {
	deptInfo, _, err := DingClient.GetDepartment(ctx, &RequestDepartmentInfo{DeptId: "155387235"})
	assert.Nil(t, err)
	assert.NotNil(t, deptInfo)
	deptInfo, _, err = DingClient.GetDepartment(ctx, &RequestDepartmentInfo{DeptId: "155387235", Language: EN_US})
	assert.Nil(t, err)
	assert.NotNil(t, deptInfo)
	_, _, err = DingClient.GetDepartment(ctx, &RequestDepartmentInfo{DeptId: ""})
	assert.NotNil(t, err)
}

func TestClient_GetUserInfoByMobileV2(t *testing.T) {
	userid, _, err := DingClient.GetUserInfoByMobileV2(ctx, testUser.Result.Mobile)
	assert.Nil(t, err)
	assert.Equal(t, testUser.Result.UserID, userid)
}

func TestClient_GetUserInfoV2(t *testing.T) {
	user, _, err := DingClient.GetUserInfoV2(ctx, &RequestUserGet{UserID: testUser.Result.UserID})
	assert.Nil(t, err)
	assert.Equal(t, user.Result.UserID, testUser.Result.UserID)
}

func TestClient_ListParentDeptByUserV2(t *testing.T) {
	_, _, err := DingClient.ListParentDeptByUserV2(ctx, testUser.Result.UserID)
	assert.Nil(t, err)
}

func TestClient_GetOrganizationUserCount(t *testing.T) {
	_, _, err := DingClient.GetOrganizationUserCount(ctx, 1)
	assert.Nil(t, err)
}

func TestClient_CreateProcessInstance(t *testing.T) {
	req := &RequestCreateProcessInstance{
		FormComponentValues: []*FormComponentValue{
			{
				Name:  "补卡时间",
				Value: "2021-03-31 09:00",
			},
			{
				Name:  "补卡理由",
				Value: "忘打卡",
			},
		},
		AgentID:          AgentID,
		DeptID:           "477274342",
		ProcessCode:      "PROC-0F200284-842E-46FF-9ACC-64DB3176150C",
		OriginatorUserID: UserID,
		ApproversV2: []ProcessApprovers{
			{
				TaskActionType: "OR",
				UserIDs:        []string{UserID, "1824661867777911"},
			},
			{
				TaskActionType: "NONE",
				UserIDs:        []string{UserID},
			},
		},
	}
	id, _, _ := DingClient.CreateProcessInstance(ctx, req)
	fmt.Println(id)
}

func TestClient_GetProcessInstance(t *testing.T) {
	//"f4b924d1-b13b-4dc1-ae5a-aea4716ecb5c"
	//"01b3a55b-d92d-40c7-ae70-87b13daf11e5"
	pi, _, _ := DingClient.GetProcessInstance(ctx, "01b3a55b-d92d-40c7-ae70-87b13daf11e5")
	b, _ := json.Marshal(pi)
	fmt.Println(string(b))
}

func TestClient_1(t *testing.T) {
	// 获取用户父部门列表
	res, _, _ := DingClient.ListParentDeptByUserV2(ctx, UserID)
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
	for _, dept := range res.ParentDeptList[0].ParentDeptIdList {
		// 获取部门详情
		resp, _, _ := DingClient.GetDepartmentV2(ctx, dept)
		b, _ = json.Marshal(resp)
		fmt.Println(string(b))
	}
}

func TestClient_GetParentDepartmentV2(t *testing.T) {
	info, _, err := DingClient.GetParentDepartmentV2(ctx, DeptID)
	assert.Nil(t, err)
	fmt.Println(info.ParentIDList)
}

func TestClient_GetSubDepartmentV2(t *testing.T) {
	info, _, err := DingClient.GetSubDepartmentV2(ctx, DeptID)
	assert.Nil(t, err)
	fmt.Println(info.SubIDList)
}
