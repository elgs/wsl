package interceptors

import (
	"database/sql"
	"fmt"

	"github.com/elgs/gostrgen"
	"github.com/elgs/wsl"
)

type ForgetPasswordSendCodeInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ForgetPasswordSendCodeInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	forgetPasswordCode, err := gostrgen.RandGen(6, gostrgen.Digit, "", "")
	if err != nil {
		return err
	}
	params := context["params"].(map[string]interface{})
	params["__forget_password"] = forgetPasswordCode
	context["forget_password"] = forgetPasswordCode

	return nil
}

func (this *ForgetPasswordSendCodeInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, scriptLabel string, cumulativeResults map[interface{}]interface{}) (bool, error) {
	if scriptIndex == 5 {
		fmt.Println(cumulativeResults)
	}

	return false, nil
}

func (this *ForgetPasswordSendCodeInterceptor) AfterEach(tx *sql.Tx, context map[string]interface{}, scriptIndex int, scriptLabel string, result interface{}, cumulativeResults map[interface{}]interface{}) error {

	if val, ok := result.([]map[string]string); ok && len(val) == 1 {
		email := val[0]["email"]
		userFlagCode := context["forget_password"]

		if wslApp, ok := context["app"].(*wsl.WSL); ok {
			err := wslApp.SendMail(email, "Password Reset Verification Code", userFlagCode.(string), email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
