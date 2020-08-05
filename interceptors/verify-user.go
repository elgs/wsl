package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type VerifyUserInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *VerifyUserInterceptor) After(tx *sql.Tx, context map[string]interface{}, results *interface{}, allResult interface{}) error {
	if sessionId, ok := context["session_id"].(string); ok {
		if flags, ok := sessions[sessionId]["flags"].(map[string]interface{}); ok {
			delete(flags, "signup")
		}
	}
	return nil
}
