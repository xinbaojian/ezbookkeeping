package models

import "time"

// CustomerType represents customer type
type CustomerType int8

// Customer types
const (
	CUSTOMER_TYPE_CUSTOMER    CustomerType = 1 // 客户
	CUSTOMER_TYPE_SUPPLIER    CustomerType = 2 // 供应商
	CUSTOMER_TYPE_BOTH        CustomerType = 3 // 既是客户也是供应商
)

type Customer struct {
	Id  int64  `xorm:"PK"`
	Uid int64  `xorm:"INDEX(IDX_customer_uid_deleted) NOT NULL"`
	Deleted bool `xorm:"INDEX(IDX_customer_uid_deleted) NOT NULL"`
	// 客户名称
	Name string `xorm:"VARCHAR(256) NOT NULL"`
	// 客户类型 1:客户 2:供应商 3:既是客户也是供应商
	CustomerType int8 `xorm:"DEFAULT 1 NOT NULL"`
	// 地址
	Address string `xorm:"VARCHAR(512) DEFAULT ''"`
	// 联系人
	Contacts string `xorm:"VARCHAR(256) DEFAULT ''"`
	// 联系方式
	ContactsInfo string `xorm:"VARCHAR(256) DEFAULT ''"`
	// 备注
	Comment string `xorm:"VARCHAR(512) DEFAULT ''"`
	// 是否隐藏
	Hidden bool `xorm:"DEFAULT false NOT NULL"`
	// 创建/更新时间
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// ToCustomerInfoResponse returns a view-object according to database model
func (c *Customer) ToCustomerInfoResponse(userTimezone *time.Location) *CustomerInfo {
	return &CustomerInfo{
		Id:           c.Id,
		Name:         c.Name,
		CustomerType: c.CustomerType,
		Address:      c.Address,
		Contacts:     c.Contacts,
		ContactsInfo: c.ContactsInfo,
		Comment:      c.Comment,
		Hidden:       c.Hidden,
		CreatedTime:  time.Unix(c.CreatedUnixTime, 0).In(userTimezone),
		UpdatedTime:  time.Unix(c.UpdatedUnixTime, 0).In(userTimezone),
	}
}