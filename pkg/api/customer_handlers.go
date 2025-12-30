package api

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// CustomersApi represents customer api
type CustomersApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	customers *services.CustomerService
}

// Initialize a customer api singleton instance
var (
	Customers = &CustomersApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		customers: services.Customers,
	}
)

// CustomerListHandler returns customers list of current user
func (a *CustomersApi) CustomerListHandler(c *core.WebContext) (any, *errs.Error) {
	var customerListReq models.CustomerListRequest
	err := c.ShouldBindQuery(&customerListReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	customers, err := a.customers.GetAllCustomersByUser(c, uid, customerListReq.CustomerType, customerListReq.VisibleOnly)

	if err != nil {
		log.Errorf(c, "[customers.CustomerListHandler] failed to get all customers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	userAllCustomerResps := make([]*models.CustomerInfo, len(customers))
	clientTimezone, _, _ := c.GetClientTimezone()

	for i := 0; i < len(customers); i++ {
		userAllCustomerResps[i] = customers[i].ToCustomerInfoResponse(clientTimezone)
	}

	return userAllCustomerResps, nil
}

// CustomerListWithPaginationHandler returns customers list with pagination of current user
func (a *CustomersApi) CustomerListWithPaginationHandler(c *core.WebContext) (any, *errs.Error) {
	var customerListReq models.CustomerListRequest
	err := c.ShouldBindQuery(&customerListReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerListWithPaginationHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	totalCount, err := a.customers.GetTotalCustomerCountByUser(c, uid, customerListReq.CustomerType, customerListReq.VisibleOnly)

	if err != nil {
		log.Errorf(c, "[customers.CustomerListWithPaginationHandler] failed to get total customer count for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	customers, err := a.customers.GetAllCustomersByUserWithPagination(c, uid, customerListReq.CustomerType, customerListReq.VisibleOnly, customerListReq.Page, customerListReq.PageSize)

	if err != nil {
		log.Errorf(c, "[customers.CustomerListWithPaginationHandler] failed to get customers for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	customerInfoResponses := make([]*models.CustomerInfo, len(customers))
	clientTimezone, _, _ := c.GetClientTimezone()

	for i, customer := range customers {
		customerInfoResponses[i] = customer.ToCustomerInfoResponse(clientTimezone)
	}

	totalPages := int(totalCount) / customerListReq.PageSize
	if int(totalCount)%customerListReq.PageSize > 0 {
		totalPages++
	}

	response := &models.CustomerListResponse{
		Total:      totalCount,
		Page:       customerListReq.Page,
		PageSize:   customerListReq.PageSize,
		TotalPages: totalPages,
		Customers:  customerInfoResponses,
	}

	return response, nil
}

// CustomerGetHandler returns one specific customer of current user
func (a *CustomersApi) CustomerGetHandler(c *core.WebContext) (any, *errs.Error) {
	var customerGetReq struct {
		Id int64 `form:"id,string" binding:"required,min=1"`
	}
	err := c.ShouldBindQuery(&customerGetReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	customer, err := a.customers.GetCustomerById(c, uid, customerGetReq.Id)

	if err != nil {
		log.Errorf(c, "[customers.CustomerGetHandler] failed to get customer \"id:%d\" for user \"uid:%d\", because %s", customerGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	clientTimezone, _, _ := c.GetClientTimezone()
	return customer.ToCustomerInfoResponse(clientTimezone), nil
}

// CustomerCreateHandler saves a new customer by request parameters for current user
func (a *CustomersApi) CustomerCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var customerCreateReq models.CustomerCreateRequest
	err := c.ShouldBindJSON(&customerCreateReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && customerCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CUSTOMER, uid, customerCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[customers.CustomerCreateHandler] another customer \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			customerId, err := utils.StringToInt64(remark)

			if err == nil {
				customer, err := a.customers.GetCustomerById(c, uid, customerId)

				if err != nil {
					log.Errorf(c, "[customers.CustomerCreateHandler] failed to get existed customer \"id:%d\" for user \"uid:%d\", because %s", customerId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				clientTimezone, _, _ := c.GetClientTimezone()
				return customer.ToCustomerInfoResponse(clientTimezone), nil
			}
		}
	}

	customer, err := a.customers.CreateCustomer(c, uid, &customerCreateReq)

	if err != nil {
		log.Errorf(c, "[customers.CustomerCreateHandler] failed to create customer for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.CustomerCreateHandler] user \"uid:%d\" has created a new customer \"id:%d\" successfully", uid, customer.Id)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CUSTOMER, uid, customerCreateReq.ClientSessionId, utils.Int64ToString(customer.Id))

	clientTimezone, _, _ := c.GetClientTimezone()
	return customer.ToCustomerInfoResponse(clientTimezone), nil
}

// CustomerModifyHandler saves an existed customer by request parameters for current user
func (a *CustomersApi) CustomerModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var customerModifyReq models.CustomerModifyRequest
	err := c.ShouldBindJSON(&customerModifyReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if customerModifyReq.Id <= 0 {
		return nil, errs.ErrCustomerIdInvalid
	}

	uid := c.GetCurrentUid()
	customer, err := a.customers.ModifyCustomer(c, uid, &customerModifyReq)

	if err != nil {
		log.Errorf(c, "[customers.CustomerModifyHandler] failed to update customer \"id:%d\" for user \"uid:%d\", because %s", customerModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.CustomerModifyHandler] user \"uid:%d\" has updated customer \"id:%d\" successfully", uid, customer.Id)

	clientTimezone, _, _ := c.GetClientTimezone()
	return customer.ToCustomerInfoResponse(clientTimezone), nil
}

// CustomerDeleteHandler deletes an existed customer by request parameters for current user
func (a *CustomersApi) CustomerDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var customerDeleteReq models.CustomerDeleteRequest
	err := c.ShouldBindJSON(&customerDeleteReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.customers.DeleteCustomer(c, uid, customerDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[customers.CustomerDeleteHandler] failed to delete customer \"id:%d\" for user \"uid:%d\", because %s", customerDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.CustomerDeleteHandler] user \"uid:%d\" has deleted customer \"id:%d\"", uid, customerDeleteReq.Id)
	return true, nil
}

// CustomerHideHandler hides an existed customer by request parameters for current user
func (a *CustomersApi) CustomerHideHandler(c *core.WebContext) (any, *errs.Error) {
	var customerHideReq models.CustomerHideRequest
	err := c.ShouldBindJSON(&customerHideReq)

	if err != nil {
		log.Warnf(c, "[customers.CustomerHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.customers.UpdateCustomerHiddenStatus(c, uid, customerHideReq.Id, customerHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[customers.CustomerHideHandler] failed to hide customer \"id:%d\" for user \"uid:%d\", because %s", customerHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[customers.CustomerHideHandler] user \"uid:%d\" has hidden customer \"id:%d\"", uid, customerHideReq.Id)
	return true, nil
}
