package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type ChangePasswordInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ChangePasswordInterceptor) Before(tx *sql.Tx, context map[string]any) error {

	if session, ok := context["session"].(map[string]any); ok {
		if flags, ok := session["flags"].(map[string]any); ok {
			if flags["signup"] != nil {
				return errors.New("user_not_verified")
			}
		}
	} else {
		return errors.New("invalid_session")
	}

	return nil
}

func (this *ChangePasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]any, statement *wsl.Statement, cumulativeResults map[string]any) (bool, error) {
	if statement.Label == "delete_sessions" {
		if cumulativeResults["change_password"] == int64(0) {
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}
