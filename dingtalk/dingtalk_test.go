package dingtalk

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var DingClient *Client
var ctx = context.Background()
var testUser *ResponseGetUserInfo

var (
	AgentID = os.Getenv("AgentID")
	AppKey = os.Getenv("AppKey")
	AppSecret = os.Getenv("AppSecret")
	UserID = os.Getenv("UserID")
)

func init()  {
	DingClient = NewClient(Option{AgentID: AgentID, AppKey: AppKey, AppSecret: AppSecret})
	var err error
	fmt.Println(len(AgentID), len(AppKey), len(AppSecret), len(UserID))
	testUser, _ , err = DingClient.GetUserInfoV2(ctx, &RequestUserGet{UserID: UserID})
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
	user, _ , err := DingClient.GetUserInfoV2(ctx, &RequestUserGet{UserID: testUser.Result.UserID})
	assert.Nil(t, err)
	assert.Equal(t, user.Result.UserID, testUser.Result.UserID)
}

func TestClient_ListParentDeptByUserV2(t *testing.T) {
	_, _ , err := DingClient.ListParentDeptByUserV2(ctx, testUser.Result.UserID)
	assert.Nil(t, err)
}

func TestClient_GetOrganizationUserCount(t *testing.T) {
	_, _, err := DingClient.GetOrganizationUserCount(ctx, 1)
	assert.Nil(t, err)
}


