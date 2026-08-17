package constants

const (
	// Validation Errors (HTTP 400)
	ErrInvalidInput                  = "Invalid input data."
	ErrMissingFieldID                = "Missing required field: ID."
	ErrMissingFieldName              = "Missing required field: Name."
	ErrMissingFieldEmail             = "Missing required field: Email."
	ErrInvalidEmailFormat            = "Invalid email format."
	ErrInvalidJSONFormat             = "Invalid JSON format."
	ErrInvalidIDFormat               = "Invalid ID format (UUID required)."
	ErrInvalidOrderType              = "Invalid order type."
	ErrInvalidFormatOrIncompleteData = "Invalid format or missing data."

	// Resource Errors (HTTP 404 / 409)
	ErrCustomerNotFound            = "Customer not found."
	ErrEmailAlreadyExists          = "The provided email is already in use."
	ErrCustomerIDEmpty             = "Customer ID cannot be empty."
	ErrCustomerNameEmpty           = "Customer name cannot be empty."
	ErrCustomerEmailEmpty          = "Customer email cannot be empty."
	ErrOrderIDEmpty                = "Order ID cannot be empty."
	ErrProductIDEmpty              = "Product ID cannot be empty."
	ErrQuantityGreaterThanZero     = "Quantity must be greater than zero."
	ErrPricePerUnitGreaterThanZero = "Price per unit must be greater than zero."
	ErrMaterialNameEmpty           = "Material name cannot be empty."
	ErrCostPerKGAtLeastZero        = "Cost per kg must be at least zero."
	ErrResourceNotFound            = "The requested resource could not be found."

	// Server Errors (HTTP 500)
	ErrInternalServer = "An unexpected error occurred. Please try again later."
)
