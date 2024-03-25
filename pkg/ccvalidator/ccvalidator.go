package ccvalidator

type Validator interface {
	Validate(number string, cvv int, month int, year int) error
}
