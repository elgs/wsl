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
		sessions[sessionId]["user_flag"] = ""
	}
	return nil
}
