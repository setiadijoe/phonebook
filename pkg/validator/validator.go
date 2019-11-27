package validator

import (
	"sync"
)

//ValidationFunc function used in server validator
type ValidationFunc func(req interface{}) error

var onceVal sync.Once
var val *Validator

//DefaultValidator execute ValidateStruct
func DefaultValidator() ValidationFunc {
	onceVal.Do(func() {
		val = New()
	})

	return func(req interface{}) error {
		return val.Validate.Struct(req)
	}
}
