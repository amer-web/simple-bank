package helper

import (
	"github.com/go-playground/validator/v10"
)

var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"EGP": true,
	// Add other valid currencies here
}

func CurrencyValidate(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if ok {
		_, exists := validCurrencies[currency]
		return exists
	}
	return false
}
