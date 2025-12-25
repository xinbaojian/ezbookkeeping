package models

type Customer struct {
	Id int64 `xorm:"PK"`
	// 客户名称
	Name    string `xorm:"VARCHAR(256) NOT NULL"`
	Uid     int64  `xorm:"INDEX(IDX_customer_uid_deleted) NOT NULL"`
	Deleted bool   `xorm:"INDEX(IDX_customer_uid_deleted) NOT NULL"`
	// 供应商
	Supplier bool `xorm:"VARCHAR(256)"`
	// 客户
	Customer bool `xorm:"VARCHAR(256)"`
	// 地址
	Address string `xorm:"VARCHAR(256) NOT NULL"`
	// 联系人
	Contacts string `xorm:"-"`
	// 联系方式
	ContactsInfo string `xorm:"-"`
	// 备注
	Comment         string `xorm:"VARCHAR(512) NOT NULL"`
	Hidden          bool   `xorm:"NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}
