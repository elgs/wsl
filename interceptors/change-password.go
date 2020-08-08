package interceptors

import (
	"database/sql"
	"errors"
	"fmt"

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

func (this *ChangePasswordInterceptor) BeforeEach(tx *sql.Tx, context map[string]interface{}, script *string, sqlParams []interface{}, scriptIndex int, scriptLabel string, cumulativeResults map[interface{}]interface{}) (bool, error) {
	if scriptLabel == "#delete_session" {
		fmt.Println("aaaa")
		if cumulativeResults["#update_password"] == int64(0) {
			fmt.Println("bbbb")
			// if password is not changed, skip deleting other sessions
			return true, nil
		}
	}
	return false, nil
}
