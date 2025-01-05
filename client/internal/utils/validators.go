package utils

import (
	gpvalidator "github.com/go-playground/validator/v10"
)

var validator *gpvalidator.Validate

func init() {
	validator = gpvalidator.New()
	validator.RegisterValidation("card_number", func(fl gpvalidator.FieldLevel) bool {
		cardNumber := fl.Field().String()
		if len(cardNumber) < 13 || len(cardNumber) > 19 {
			return false
		}
		return luhnCheck(cardNumber)
	})
}

func luhnCheck(cardNumber string) bool {
	var sum int
	alt := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		n := int(cardNumber[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

func StringIsValid(v string) bool {
	return v != ""
}

func CardNumberIsValid(v string) bool {
	return validator.Var(v, "card_number") == nil
}

func ExpireDateIsValid(v string) bool {
	return validator.Var(v, "required,datetime=01/06") == nil
}

func CSCIsValid(v string) bool {
	return validator.Var(v, "required,len=3") == nil
}
