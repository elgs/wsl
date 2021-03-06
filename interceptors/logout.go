package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type LogoutInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *LogoutInterceptor) After(tx *sql.Tx, context map[string]interface{}, results *interface{}, allResults interface{}) error {
	// params := context["params"].(map[string]interface{})
	if val, ok := context["session_id"].(string); ok {
		delete(Sessions, val)
	}
	return nil
}
