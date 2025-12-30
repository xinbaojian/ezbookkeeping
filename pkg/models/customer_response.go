package models

import "time"

// CustomerInfo 客户信息响应
type CustomerInfo struct {
	// 客户ID
	Id int64 `json:"id,string"`
	// 客户名称
	Name string `json:"name"`
	// 客户类型 1:客户 2:供应商 3:两者
	CustomerType int8 `json:"customer_type"`
	// 地址
	Address string `json:"address"`
	// 联系人
	Contacts string `json:"contacts"`
	// 联系方式
	ContactsInfo string `json:"contacts_info"`
	// 备注
	Comment string `json:"comment"`
	// 是否隐藏
	Hidden bool `json:"hidden"`
	// 创建时间
	CreatedTime time.Time `json:"created_time"`
	// 更新时间
	UpdatedTime time.Time `json:"updated_time"`
}

// CustomerListResponse 客户列表响应
type CustomerListResponse struct {
	// 总数量
	Total int64 `json:"total"`
	// 当前页
	Page int `json:"page"`
	// 分页大小
	PageSize int `json:"page_size"`
	// 总页数
	TotalPages int `json:"total_pages"`
	// 数据列表
	Customers []*CustomerInfo `json:"customers"`
}

// CustomerGetResponse 获取单个客户响应
type CustomerGetResponse struct {
	CustomerInfo
}

// CustomerCreateResponse 创建客户响应
type CustomerCreateResponse struct {
	CustomerInfo
}

// CustomerModifyResponse 修改客户响应
type CustomerModifyResponse struct {
	CustomerInfo
}