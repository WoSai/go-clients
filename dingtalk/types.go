package dingtalk

import (
	"encoding/json"
	"time"
)

type (
	Response interface {
		GotErr() error
	}

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

	UserGetResponse struct {
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
		Result        *UserGetResponse `json:"result"`
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

	DeptLeader struct {
		DeptId int  `json:"dept_id"`
		Leader bool `json:"leader"`
	}

	DeptOrder struct {
		DeptId int `json:"dept_id"`
		Order  int `json:"order"`
	}

	CompleteUserInfo struct {
		UserID           string       `json:"userid"`
		UnionID          string       `json:"unionid"`
		Name             string       `json:"name"`
		Avatar           string       `json:"avatar"`
		StateCode        string       `json:"state_code"`
		ManagerUserID    string       `json:"manager_userid"`
		Mobile           string       `json:"mobile"`
		HideMobile       bool         `json:"hide_mobile"`
		Telephone        string       `json:"telephone"`
		JobNumber        string       `json:"job_number"`
		Title            string       `json:"title"`
		Email            string       `json:"email"`
		WorkPlace        string       `json:"work_place"`
		Remark           string       `json:"remark"`
		ExclusiveAccount bool         `json:"exclusive_account"`
		OrgEmail         string       `json:"org_email"`
		DeptIdList       []int        `json:"dept_id_list"`
		DeptOrderList    []DeptOrder  `json:"dept_order_list"`
		Extension        string       `json:"extension"`
		HiredDate        int          `json:"hired_date"`
		Active           bool         `json:"active"`
		RealAuthed       bool         `json:"real_authed"`
		Senior           bool         `json:"senior"`
		Admin            bool         `json:"admin"`
		Boss             bool         `json:"boss"`
		LeaderInDept     []DeptLeader `json:"leader_in_dept"`
		RoleList         []UserRole   `json:"role_list"`
	}

	ResponseGetUserInfo struct {
		*DingtalkErr `json:",inline"`
		Result       *CompleteUserInfo `json:"result"`
	}

	ResponseUserByMobile struct {
		BasicResponse `json:",inline"`
		Result        struct {
			Userid string `json:"userid"`
		} `json:"result"`
	}

	ParentDeptIds struct {
		ParentDeptIdList []int `json:"parent_dept_id_list"`
	}

	DeptListParent struct {
		ParentDeptList []ParentDeptIds `json:"parent_list"`
	}

	ResponseGetUserIdByUnionid struct {
		BasicResponse `json:",inline"`
		Result        *DeptListParent `json:"result"`
	}

	RequestDepartmentInfo struct {
		DeptId   string `json:"id"`
		Language Lang   `json:"lang,omitempty"`
	}

	DepartmentInfo struct {
		Id                    int    `json:"id"`
		Name                  string `json:"name"`
		Order                 int    `json:"order"`
		ParentId              int    `json:"parentid"`
		SourceIdentifier      string `json:"source_identifier"`
		CreateDeptGroup       bool   `json:"createDeptGroup"`
		AutoAddUser           bool   `json:"autoAddUser"`
		GroupContainSubDept   bool   `json:"groupContainSubDept"`
		OrgDeptOwner          string `json:"orgDeptOwner"`
		DeptGroupChatId       string `json:"deptGroupChatId"`
		DeptManagerUseridList string `json:"deptManagerUseridList"`
		OuterDept             bool   `json:"outerDept"`
		OuterPermitUsers      string `json:"outerPermitUsers"`
		OuterPermitDepts      string `json:"outerPermitDepts"`
		DeptHiding            bool   `json:"deptHiding"`
		DeptPermits           string `json:"deptPermits"`
		UserPermits           string `json:"userPermits"`
	}

	ResponseDeptInfo struct {
		BasicResponse `json:",inline"`
		*DepartmentInfo
	}

	RequestDepartmentInfoV2 struct {
		DeptID   int    `json:"dept_id"`
		Language string `json:"language,omitempty"`
	}

	RequestDepartmentUserList struct {
		DeptID int `json:"dept_id"`
		Cursor int `json:"cursor"`
		Size   int `json:"size"`
	}

	DepartmentUserInfo struct {
		UserID               string `json:"userid"`
		UnionID              string `json:"unionid"`
		Name                 string `json:"name"`
		Avatar               string `json:"avatar"`
		StateCode            string `json:"state_code"`
		Mobile               string `json:"mobile"`
		HideMobile           bool   `json:"hide_mobile"`
		Telephone            string `json:"telephone"`
		JobNumber            string `json:"job_number"`
		Title                string `json:"title"`
		Email                string `json:"email"`
		OrgEmail             string `json:"org_email"`
		WorkPlace            string `json:"work_place"`
		Remark               string `json:"remark"`
		DeptIDList           []int  `json:"dept_id_list"`
		DeptOrder            int    `json:"dept_order"`
		Extension            string `json:"extension"`
		HiredDate            int    `json:"hired_date"`
		Active               bool   `json:"active"`
		Admin                bool   `json:"admin"`
		Boss                 bool   `json:"boss"`
		Leader               bool   `json:"leader"`
		ExclusiveAccount     bool   `json:"exclusive_account"`
		LoginID              string `json:"login_id"`
		ExclusiveAccountType string `json:"exclusive_account_type"`
	}

	DepartmentUserList struct {
		HasMore    bool                 `json:"has_more"`
		NextCursor int                  `json:"next_cursor"`
		List       []DepartmentUserInfo `json:"list"`
	}

	RespDepartmentUserInfo struct {
		BasicResponse `json:",inline"`
		Result        *DepartmentUserList `json:"result,omitempty"`
	}

	RequestDepartmentListParentByUser struct {
		UserID string `json:"userid"`
	}

	DepartmentListParentByUser struct {
		ParentList []struct {
			ParentDeptIDList []int `json:"parent_dept_id_list"`
		} `json:"parent_list"`
	}

	RespDepartmentListParentByUser struct {
		BasicResponse `json:",inline"`
		Result        *DepartmentListParentByUser `json:"result,omitempty"`
	}

	DepartmentInfoV2 struct {
		DeptID                int      `json:"dept_id"`
		Name                  string   `json:"name,omitempty"`
		ParentID              int      `json:"parent_id,omitempty"`
		SourceIdentifier      string   `json:"source_identifier,omitempty"`
		CreateDeptGroup       bool     `json:"create_dept_group,omitempty"`
		AutoAddUser           bool     `json:"auto_add_user,omitempty"`
		FromUnionOrg          bool     `json:"from_union_org,omitempty"`
		Tags                  string   `json:"tags,omitempty"`
		Order                 int      `json:"order,omitempty"`
		DeptGroupChatID       string   `json:"dept_group_chat_id,omitempty"`
		GroupContainSubDept   bool     `json:"group_contain_sub_dept,omitempty"`
		OrgDeptOwner          string   `json:"org_dept_owner,omitempty"`
		DeptManagerUseridList []string `json:"dept_manager_userid_list,omitempty"`
		OuterDept             bool     `json:"outer_dept,omitempty"`
		OuterPermitDepts      []int    `json:"outer_permit_depts,omitempty"`
		OuterPermitUsers      []string `json:"outer_permit_users,omitempty"`
		HideDept              bool     `json:"hide_dept,omitempty"`
		UserPermits           []string `json:"user_permits,omitempty"`
		DeptPermits           []int    `json:"dept_permits,omitempty"`
	}

	ResponseDeptInfoV2 struct {
		BasicResponse `json:",inline"`
		Result        *DepartmentInfoV2 `json:"result,omitempty"`
	}

	ParentDeptsV2 struct {
		ParentIDList []int `json:"parent_id_list"`
	}

	RespParentDeptsV2 struct {
		BasicResponse `json:",inline"`
		Result        *ParentDeptsV2 `json:"result,omitempty"`
	}

	SubDeptsV2 struct {
		SubIDList []int `json:"dept_id_list"`
	}

	RespSubDeptsV2 struct {
		BasicResponse `json:",inline"`
		Result        *SubDeptsV2 `json:"result,omitempty"`
	}

	// 钉钉的时间格式
	DingTime struct {
		time.Time
	}

	FormComponentValue struct {
		Name          string `json:"name,omitempty"`
		Value         string `json:"value,omitempty"`
		ExtValue      string `json:"ext_value,omitempty"`
		ID            string `json:"id,omitempty"`
		ComponentType string `json:"component_type,omitempty"`
	}

	ProcessApprovers struct {
		TaskActionType string   `json:"task_action_type,omitempty"`
		UserIDs        []string `json:"user_ids,omitempty"`
	}

	// 审批请求参数
	RequestCreateProcessInstance struct {
		FormComponentValues []*FormComponentValue `json:"form_component_values,omitempty"`
		AgentID             string                `json:"agent_id,omitempty"`
		DeptID              string                `json:"dept_id,omitempty"`
		ProcessCode         string                `json:"process_code,omitempty"`
		OriginatorUserID    string                `json:"originator_user_id,omitempty"` // 审批发起人ID
		ApproversV2         []ProcessApprovers    `json:"approvers_v2,omitempty"`
		CCList              string                `json:"cc_list,omitempty"`     // 抄送人userid列表
		CCPosition          string                `json:"cc_position,omitempty"` // 抄送时间，分为（START, FINISH, START_FINISH）
	}

	ResponseCreateProcessInstance struct {
		BasicResponse     `json:",inline"`
		ProcessInstanceID string `json:"process_instance_id,omitempty"`
	}

	ProcessAttachment struct {
		FileName string `json:"file_name,omitempty"`
		FileSize string `json:"file_size,omitempty"`
		FileID   string `json:"file_id,omitempty"`
		FileType string `json:"file_type,omitempty"`
	}

	ProcessOperationRecord struct {
		UserID          string               `json:"user_id,omitempty"`
		Date            DingTime             `json:"date,omitempty"`
		OperationType   string               `json:"operation_type,omitempty"`
		OperationResult string               `json:"operation_result,omitempty"`
		Remark          string               `json:"remark,omitempty"`
		Attachments     []*ProcessAttachment `json:"attachments,omitempty"`
	}

	ProcessTask struct {
		UserID     string   `json:"userid,omitempty"`
		TaskStatus string   `json:"task_status,omitempty"`
		TaskResult string   `json:"task_result,omitempty"`
		CreateTime DingTime `json:"create_time,omitempty"`
		FinishTime DingTime `json:"finish_time,omitempty"`
		TaskID     string   `json:"taskid,omitempty"`
		URL        string   `json:"url,omitempty"`
	}

	// 审批实例详情
	ProcessInstance struct {
		Title                      string                    `json:"title,omitempty"`
		CreateTime                 DingTime                  `json:"create_time,omitempty"`
		FinishTime                 DingTime                  `json:"finish_time,omitempty"`
		OriginatorUserID           string                    `json:"originator_userid,omitempty"`
		OriginatorDeptID           string                    `json:"originator_dept_id,omitempty"`
		Status                     string                    `json:"status,omitempty"` // RUNNING审批中COMPLETED完成
		ApproverUserIDs            []string                  `json:"approver_userids,omitempty"`
		CCUserIDs                  []string                  `json:"cc_userids,omitempty"`
		Result                     string                    `json:"result,omitempty"` // agree同意refuse拒绝
		BusinessID                 string                    `json:"business_id,omitempty"`
		OperationRecords           []*ProcessOperationRecord `json:"operation_records,omitempty"`
		Tasks                      []*ProcessTask            `json:"tasks,omitempty"`
		OriginatorDeptName         string                    `json:"originator_dept_name,omitempty"`
		BizAction                  string                    `json:"biz_action,omitempty"`
		AttachedProcessInstanceIDs []string                  `json:"attached_process_instance_ids,omitempty"`
		FormComponentValues        []*FormComponentValue     `json:"form_component_values,omitempty"`
		MainProcessInstanceID      string                    `json:"main_process_instance_id,omitempty"`
	}

	ResponseGetProcessInstance struct {
		BasicResponse   `json:",inline"`
		ProcessInstance *ProcessInstance `json:"process_instance,omitempty"`
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

func (dt *DingTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	ti, err := time.Parse("\"2006-01-02 15:04:05\"", string(b))
	if err != nil {
		return err
	}
	dt.Time = ti
	return nil
}

func (dt DingTime) MarshalJSON() ([]byte, error) {
	return []byte(dt.Time.Format("\"2006-01-02 15:04:05\"")), nil
}
