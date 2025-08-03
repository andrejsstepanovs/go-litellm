package users

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID int64 `json:"id"`
}

func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("user is required")
	}
	var validate = validator.New()
	err := validate.Struct(u)
	if err != nil {
		return fmt.Errorf("user validation error: %w", err)
	}
	return nil
}
