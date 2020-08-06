package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type ResetPasswordInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ResetPasswordInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	if context["session_id"] == "" {
		return errors.New("invalid_token")
	}

	if session, ok := context["session"].(map[string]string); ok {
		if session["user_flag"] == "signup" {
			return errors.New("user_not_verified")
		}
	} else {
		return errors.New("invalid_session")
	}

	if context["user_mode"] != "root" {
		return errors.New("access_denied")
	}

	return nil
}

func (this *ResetPasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, cumulativeResults []interface{}) (bool, error) {
	if scriptIndex == 3 {
		if cumulativeResults[2] == int64(0) {
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}