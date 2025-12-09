package validator

import (
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func GetValidator() *validator.Validate {
	if validate != nil {
		return validate
	}
	// v11以降でデフォルトの動作となる新しい動作を有効にするオプションを使用して初期化
	validate = validator.New(validator.WithRequiredStructEnabled())
	return validate
}