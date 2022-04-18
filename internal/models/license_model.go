package models

type Licence struct {
	ID             int
	CustomerId     int
	ProductType    string
	CallbackUrl    string
	Count          int
	LicenseKey     string
	RegistrationAt string
	ActivationAt   string
	ExpirationAt   string
	Duration       int
	DeletedAt      string
}
