package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type ResetPasswordInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ResetPasswordInterceptor) Before(tx *sql.Tx, context map[string]any) error {

	if session, ok := context["session"].(map[string]any); ok {
		if flags, ok := session["flags"].(map[string]any); ok {
			if flags["signup"] != nil {
				return errors.New("user_not_verified")
			}
		}
	} else {
		return errors.New("invalid_session")
	}

	if context["user_mode"] != "root" {
		return errors.New("access_denied")
	}

	return nil
}

func (this *ResetPasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]any, script *string, sqlParams *[]any, scriptIndex int, scriptLabel string, cumulativeResults map[string]any) (bool, error) {
	if scriptLabel == "delete_sessions" {
		if cumulativeResults["reset_password"] == int64(0) {
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}
