package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type VerifyUserInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *VerifyUserInterceptor) After(tx *sql.Tx, context map[string]interface{}, results *interface{}, allResult interface{}) error {
	if session, ok := context["session"].(map[string]interface{}); ok {
		if flags, ok := session["flags"].(map[string]interface{}); ok {
			delete(flags, "signup")
		}
	}
	return nil
}
