package models

import "time"

type Licence struct {
	ID               int
	Code             string
	CustomerId       int
	ProductType      string
	CallbackURL      string
	IsSentCallback   bool
	CallbackAttempts uint
	Count            int
	LicenseKey       string
	RegistrationAt   time.Time
	ActivationAt     time.Time
	ExpirationAt     time.Time
	Duration         int
	Description      string
	DeletedAt        time.Time
	Customer         Customer
}
