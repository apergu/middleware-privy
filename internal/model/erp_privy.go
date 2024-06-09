package model

import "github.com/go-playground/validator/v10"

type CheckTopUpStatus struct {
	TopUPID string `json:"topup_id" validate:"required"`
	Event   string `json:"event" validate:"required"`
}

func (c CheckTopUpStatus) Validate() []map[string]interface{} {
	var validationErrors []map[string]interface{}
	v := validator.New()

	// Create a user with invalid data
	// Validate the user struct
	err := v.Struct(c)
	if err != nil {
		// Validation failed, print the error messages
		for _, err := range err.(validator.ValidationErrors) {
			//fmt.Println(err)
			//return err
			validationErrors = append(validationErrors,
				map[string]interface{}{
					"field":       err.Field(),
					"description": err.Tag(),
				})
		}
	}

	return validationErrors
}

type VoidBalance struct {
	TopUPID string `json:"topup_id" validate:"required"`
}

func (c VoidBalance) Validate() []map[string]interface{} {
	var validationErrors []map[string]interface{}
	v := validator.New()

	// Create a user with invalid data
	// Validate the user struct
	err := v.Struct(c)
	if err != nil {
		// Validation failed, print the error messages
		for _, err := range err.(validator.ValidationErrors) {
			//fmt.Println(err)
			//return err
			validationErrors = append(validationErrors,
				map[string]interface{}{
					"field":       err.Field(),
					"description": err.Tag(),
				})
		}
	}

	return validationErrors
}
