package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type ForgetPasswordResetPasswordInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ForgetPasswordResetPasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, scriptLabel string, cumulativeResults map[string]interface{}) (bool, error) {
	if scriptLabel == "delete_session" {
		if cumulativeResults["delete_flag"] == int64(0) {
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}
