package email

import (
	"fmt"

	"github.com/badoux/checkmail"
)

func ValidateEmail(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return fmt.Errorf("email is invalid")
	}
	return nil
}
