package errs

import "net/http"

// Error codes related to customers
var (
	ErrCustomerIdInvalid         = NewNormalError(NormalSubcategoryCustomer, 0, http.StatusBadRequest, "customer id is invalid")
	ErrCustomerNotFound          = NewNormalError(NormalSubcategoryCustomer, 1, http.StatusBadRequest, "customer not found")
	ErrCustomerTypeInvalid       = NewNormalError(NormalSubcategoryCustomer, 2, http.StatusBadRequest, "customer type is invalid")
	ErrCustomerNameRequired      = NewNormalError(NormalSubcategoryCustomer, 3, http.StatusBadRequest, "customer name is required")
	ErrCustomerNameTooLong       = NewNormalError(NormalSubcategoryCustomer, 4, http.StatusBadRequest, "customer name is too long")
	ErrCustomerAddressTooLong    = NewNormalError(NormalSubcategoryCustomer, 5, http.StatusBadRequest, "customer address is too long")
	ErrCustomerContactsTooLong   = NewNormalError(NormalSubcategoryCustomer, 6, http.StatusBadRequest, "customer contacts name is too long")
	ErrCustomerContactsInfoTooLong = NewNormalError(NormalSubcategoryCustomer, 7, http.StatusBadRequest, "customer contacts info is too long")
	ErrCustomerCommentTooLong    = NewNormalError(NormalSubcategoryCustomer, 8, http.StatusBadRequest, "customer comment is too long")
	ErrCustomerInUseCannotBeDeleted = NewNormalError(NormalSubcategoryCustomer, 9, http.StatusBadRequest, "customer is in use and cannot be deleted")
)