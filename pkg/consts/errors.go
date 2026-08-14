package constants

const (
	// Validation Errors (HTTP 400)
	ErrInvalidInput       = "Invalid input data."
	ErrMissingFieldID     = "Missing required field: ID."
	ErrMissingFieldName   = "Missing required field: Name."
	ErrMissingFieldEmail  = "Missing required field: Email."
	ErrInvalidEmailFormat = "Invalid email format."
	ErrInvalidJSONFormat  = "Invalid JSON format."
	ErrInvalidIDFormat    = "Invalid ID format (UUID required)."
	ErrInvalidOrderType   = "Invalid Order Type"

	// Resource Errors (HTTP 404 / 409)
	ErrCustomerNotFound            = "Customer not found."
	ErrEmailAlreadyExists          = "The provided email is already in use."
	ErrUnableUpdateData            = "Unable to update customer data."
	ErrCustomerIDEmpty             = "Customer ID name cannot be empty."
	ErrCustomerNameEmpty           = "Customer name cannot be empty."
	ErrCustomerEmailEmpty          = "Customer email cannot be empty."
	ErrOrderIDEmpty                = "Order ID cannot be empty."
	ErrProdductIDEmpty             = "Product ID cannot be empty."
	ErrQuantityGreaterThanZero     = "Quantity must be greater than zero."
	ErrPricePerUnitGreaterThanZero = "PricePerUnit must be greater than zero."
	ErrDataNotFound                = "Data not found."

	// Server Errors (HTTP 500)
	ErrInternalServer       = "An unexpected error occurred. Please try again later."
	ErrCreateInternalServer = "An unable to create data."
	ErrUpdateInternalServer = "An unable to update data."
	ErrDeleteInternalServer = "An unable to delete data."
	ErrGenerateIDServer     = "The generated ID format is incorrect."
)
