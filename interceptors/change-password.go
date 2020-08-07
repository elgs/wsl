package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type ChangePasswordInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *ChangePasswordInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	if session, ok := context["session"].(map[string]interface{}); ok {
		if flags, ok := session["flags"].(map[string]interface{}); ok {
			if flags["signup"] != nil {
				return errors.New("user_not_verified")
			}
		}
	} else {
		return errors.New("invalid_session")
	}

	return nil
}

func (this *ChangePasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, cumulativeResults []interface{}) (bool, error) {
	if scriptIndex == 4 {
		if cumulativeResults[3] == int64(0) {
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}
