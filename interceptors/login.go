package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
	"github.com/pkg/errors"
)

type LoginInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *LoginInterceptor) AfterEach(tx *sql.Tx, context map[string]interface{}, scriptIndex int, scriptLabel string, result interface{}, cumulativeResults map[interface{}]interface{}) error {

	if scriptIndex == 5 {
		if val, ok := result.([]map[string]string); ok && len(val) == 1 {
			context["session_id"] = val[0]["session_id"]
			email := val[0]["email"]
			userFlag := val[0]["code"]
			userFlagCode := val[0]["value"]

			if wslApp, ok := context["app"].(*wsl.WSL); ok && userFlag == "signup" {
				err := wslApp.SendMail(email, "New Account Verification Code", userFlagCode, email)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (this *LoginInterceptor) After(tx *sql.Tx, context map[string]interface{}, results *interface{}, allResult interface{}) error {
	if context["session_id"] == nil {
		return errors.New("login_failed")
	}
	*results = map[string]interface{}{"access_token": context["session_id"]}
	return nil
}
