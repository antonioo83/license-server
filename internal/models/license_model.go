package models

import "time"

type Licence struct {
	ID               int
	Code             string `copier:"LicenseId"`
	CustomerId       int
	ProductType      string `copier:"ProductType"`
	CallbackURL      string `copier:"CallbackURL"`
	IsSentCallback   bool
	CallbackAttempts uint
	Count            int    `copier:"Count"`
	LicenseKey       string `copier:"LicenseKey"`
	RegistrationAt   time.Time
	ActivationAt     time.Time
	ExpirationAt     time.Time
	Duration         int
	Description      string `copier:"Description"`
	DeletedAt        time.Time
	Customer         Customer
}
