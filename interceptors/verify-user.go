package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type VerifyUserInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *VerifyUserInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	if _, ok := context["session"].(map[string]interface{}); !ok {
		return errors.New("invalid_session")
	}

	return nil
}

func (this *VerifyUserInterceptor) After(tx *sql.Tx, context map[string]interface{}, results *interface{}, allResults interface{}) error {
	if session, ok := context["session"].(map[string]interface{}); ok {
		if flags, ok := session["flags"].(map[string]interface{}); ok {
			delete(flags, "signup")
		}
	}
	return nil
}
