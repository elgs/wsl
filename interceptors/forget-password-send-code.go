package interceptors

import (
	"database/sql"

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

func (this *ForgetPasswordSendCodeInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, scriptLabel string, cumulativeResults map[string]interface{}) (bool, error) {

	if skipAll, ok := context["skip_all"].(bool); ok && skipAll {
		return true, nil
	}

	return false, nil
}

func (this *ForgetPasswordSendCodeInterceptor) AfterEach(tx *sql.Tx, context map[string]interface{}, scriptIndex int, scriptLabel string, result interface{}, cumulativeResults map[string]interface{}) error {

	if scriptLabel == "get_uid_email" {
		if results, ok := result.([]map[string]string); ok && len(results) == 1 {
			email := results[0]["email"]
			userFlagCode := context["forget_password"]
			if results[0]["uid"] != "" && email != "" {
				if wslApp, ok := context["app"].(*wsl.WSL); ok {
					err := wslApp.SendMail(email, "Password Reset Verification Code", userFlagCode.(string), email)
					if err != nil {
						return err
					}
				}
			} else {
				context["skip_all"] = true
			}
		} else {
			context["skip_all"] = true
		}
	}

	return nil
}
