package models

// CustomerListRequest 客户列表请求
type CustomerListRequest struct {
	// 仅显示未隐藏的
	VisibleOnly bool `form:"visible_only"`
	// 客户类型过滤 0:全部 1:客户 2:供应商 3:两者
	CustomerType int8 `form:"customer_type,default=0"`
	// 分页参数
	Page      int `form:"page,default=1"`
	PageSize  int `form:"page_size,default=20"`
}

// CustomerCreateRequest 创建客户请求
type CustomerCreateRequest struct {
	ClientSessionId string `json:"client_session_id" binding:"omitempty"`
	// 客户名称
	Name string `json:"name" binding:"required,max=256"`
	// 客户类型
	CustomerType int8 `json:"customer_type" binding:"gte=1,lte=3"`
	// 地址
	Address string `json:"address" binding:"omitempty,max=512"`
	// 联系人
	Contacts string `json:"contacts" binding:"omitempty,max=256"`
	// 联系方式
	ContactsInfo string `json:"contacts_info" binding:"omitempty,max=256"`
	// 备注
	Comment string `json:"comment" binding:"omitempty,max=512"`
	// 是否隐藏
	Hidden bool `json:"hidden"`
}

// CustomerModifyRequest 修改客户请求
type CustomerModifyRequest struct {
	// 客户ID
	Id int64 `json:"id,string" binding:"required,gt=0"`
	// 客户名称
	Name string `json:"name" binding:"required,max=256"`
	// 客户类型
	CustomerType int8 `json:"customer_type" binding:"gte=1,lte=3"`
	// 地址
	Address string `json:"address" binding:"omitempty,max=512"`
	// 联系人
	Contacts string `json:"contacts" binding:"omitempty,max=256"`
	// 联系方式
	ContactsInfo string `json:"contacts_info" binding:"omitempty,max=256"`
	// 备注
	Comment string `json:"comment" binding:"omitempty,max=512"`
	// 是否隐藏
	Hidden bool `json:"hidden"`
}

// CustomerDeleteRequest 删除客户请求
type CustomerDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,gt=0"`
}

// CustomerHideRequest 隐藏/显示客户请求
type CustomerHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,gt=0"`
	Hidden bool  `json:"hidden"`
}