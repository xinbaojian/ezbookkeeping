package services

import (
	"time"
	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// CustomerService represents customer service
type CustomerService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a customer service singleton instance
var (
	Customers = &CustomerService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllCustomersByUser returns all customer models of user
func (s *CustomerService) GetAllCustomersByUser(c core.Context, uid int64, customerType int8, visibleOnly bool) ([]*models.Customer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var customers []*models.Customer
	session := s.UserDataDB(uid).NewSession(c)

	if customerType > 0 {
		session = session.Where("customer_type=?", customerType)
	}

	if visibleOnly {
		session = session.Where("hidden=?", false)
	}

	err := session.Where("uid=? AND deleted=?", uid, false).
		OrderBy("created_unix_time desc").
		Find(&customers)

	if err != nil {
		log.Errorf(c, "[customers.GetAllCustomersByUser] failed to get customers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return customers, nil
}

// GetTotalCustomerCountByUser returns total customer count of user
func (s *CustomerService) GetTotalCustomerCountByUser(c core.Context, uid int64, customerType int8, visibleOnly bool) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	session := s.UserDataDB(uid).NewSession(c)

	if customerType > 0 {
		session = session.Where("customer_type=?", customerType)
	}

	if visibleOnly {
		session = session.Where("hidden=?", false)
	}

	count, err := session.Where("uid=? AND deleted=?", uid, false).
		Count(&models.Customer{})

	if err != nil {
		log.Errorf(c, "[customers.GetTotalCustomerCountByUser] failed to get customer count for user \"uid:%d\", because %s", uid, err.Error())
		return 0, errs.Or(err, errs.ErrOperationFailed)
	}

	return count, nil
}

// GetAllCustomersByUserWithPagination returns customer models of user with pagination
func (s *CustomerService) GetAllCustomersByUserWithPagination(c core.Context, uid int64, customerType int8, visibleOnly bool, page, pageSize int) ([]*models.Customer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	var customers []*models.Customer
	session := s.UserDataDB(uid).NewSession(c)

	if customerType > 0 {
		session = session.Where("customer_type=?", customerType)
	}

	if visibleOnly {
		session = session.Where("hidden=?", false)
	}

	err := session.Where("uid=? AND deleted=?", uid, false).
		OrderBy("created_unix_time desc").
		Limit(pageSize, (page-1)*pageSize).
		Find(&customers)

	if err != nil {
		log.Errorf(c, "[customers.GetAllCustomersByUserWithPagination] failed to get customers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return customers, nil
}

// GetCustomerById returns customer model according to customer id
func (s *CustomerService) GetCustomerById(c core.Context, uid int64, customerId int64) (*models.Customer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if customerId <= 0 {
		return nil, errs.ErrCustomerIdInvalid
	}

	customer := &models.Customer{}
	has, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND id=?", uid, false, customerId).Get(customer)

	if err != nil {
		log.Errorf(c, "[customers.GetCustomerById] failed to get customer \"id:%d\" for user \"uid:%d\", because %s", customerId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	if !has {
		return nil, errs.ErrCustomerNotFound
	}

	return customer, nil
}

// CreateCustomer saves a new customer model to database
func (s *CustomerService) CreateCustomer(c core.Context, uid int64, customerCreateReq *models.CustomerCreateRequest) (*models.Customer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if customerCreateReq == nil {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	now := time.Now().Unix()

	customerId := s.GenerateUuid(uuid.UUID_TYPE_ACCOUNT)
	if customerId < 1 {
		return nil, errs.ErrSystemIsBusy
	}

	customer := &models.Customer{}
	customer.Id = customerId

	customer.Uid = uid
	customer.Name = customerCreateReq.Name
	customer.CustomerType = customerCreateReq.CustomerType
	customer.Address = customerCreateReq.Address
	customer.Contacts = customerCreateReq.Contacts
	customer.ContactsInfo = customerCreateReq.ContactsInfo
	customer.Comment = customerCreateReq.Comment
	customer.Hidden = customerCreateReq.Hidden
	customer.Deleted = false
	customer.CreatedUnixTime = now
	customer.UpdatedUnixTime = now
	customer.DeletedUnixTime = 0

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(customer)
		return err
	})

	if err != nil {
		log.Errorf(c, "[customers.CreateCustomer] failed to create customer for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.CreateCustomer] user \"uid:%d\" has created a new customer \"id:%d\" successfully", uid, customer.Id)

	return customer, nil
}

// ModifyCustomer saves an existed customer model to database
func (s *CustomerService) ModifyCustomer(c core.Context, uid int64, customerModifyReq *models.CustomerModifyRequest) (*models.Customer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if customerModifyReq == nil {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	customer, err := s.GetCustomerById(c, uid, customerModifyReq.Id)
	if err != nil {
		return nil, err
	}

	anythingUpdated := false
	now := time.Now().Unix()

	if customer.Name != customerModifyReq.Name {
		customer.Name = customerModifyReq.Name
		anythingUpdated = true
	}

	if customer.CustomerType != customerModifyReq.CustomerType {
		customer.CustomerType = customerModifyReq.CustomerType
		anythingUpdated = true
	}

	if customer.Address != customerModifyReq.Address {
		customer.Address = customerModifyReq.Address
		anythingUpdated = true
	}

	if customer.Contacts != customerModifyReq.Contacts {
		customer.Contacts = customerModifyReq.Contacts
		anythingUpdated = true
	}

	if customer.ContactsInfo != customerModifyReq.ContactsInfo {
		customer.ContactsInfo = customerModifyReq.ContactsInfo
		anythingUpdated = true
	}

	if customer.Comment != customerModifyReq.Comment {
		customer.Comment = customerModifyReq.Comment
		anythingUpdated = true
	}

	if customer.Hidden != customerModifyReq.Hidden {
		customer.Hidden = customerModifyReq.Hidden
		anythingUpdated = true
	}

	if !anythingUpdated {
		return nil, errs.ErrNothingWillBeUpdated
	}

	customer.UpdatedUnixTime = now

	err = s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err = sess.ID(customer.Id).Cols("name", "customer_type", "address", "contacts", "contacts_info", "comment", "hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(customer)
		return err
	})

	if err != nil {
		log.Errorf(c, "[customers.ModifyCustomer] failed to update customer \"id:%d\" for user \"uid:%d\", because %s", customerModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.ModifyCustomer] user \"uid:%d\" has updated customer \"id:%d\" successfully", uid, customer.Id)

	return customer, nil
}

// DeleteCustomer deletes an existed customer from database
func (s *CustomerService) DeleteCustomer(c core.Context, uid int64, customerId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if customerId <= 0 {
		return errs.ErrCustomerIdInvalid
	}

	customer, err := s.GetCustomerById(c, uid, customerId)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	customer.Deleted = true
	customer.DeletedUnixTime = now
	customer.UpdatedUnixTime = now

	err = s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err = sess.ID(customer.Id).Cols("deleted", "deleted_unix_time", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(customer)
		return err
	})

	if err != nil {
		log.Errorf(c, "[customers.DeleteCustomer] failed to delete customer \"id:%d\" for user \"uid:%d\", because %s", customerId, uid, err.Error())
		return errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.DeleteCustomer] user \"uid:%d\" has deleted customer \"id:%d\" successfully", uid, customerId)

	return nil
}

// UpdateCustomerHiddenStatus updates hidden field of given customer
func (s *CustomerService) UpdateCustomerHiddenStatus(c core.Context, uid int64, customerId int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if customerId <= 0 {
		return errs.ErrCustomerIdInvalid
	}

	customer, err := s.GetCustomerById(c, uid, customerId)
	if err != nil {
		return err
	}

	if customer.Hidden == hidden {
		return nil
	}

	now := time.Now().Unix()
	customer.Hidden = hidden
	customer.UpdatedUnixTime = now

	err = s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err = sess.ID(customer.Id).Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(customer)
		return err
	})

	if err != nil {
		log.Errorf(c, "[customers.UpdateCustomerHiddenStatus] failed to update hidden status of customer \"id:%d\" for user \"uid:%d\", because %s", customerId, uid, err.Error())
		return errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.UpdateCustomerHiddenStatus] user \"uid:%d\" has updated hidden status of customer \"id:%d\"", uid, customerId)

	return nil
}
