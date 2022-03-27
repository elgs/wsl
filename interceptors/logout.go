package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type LogoutInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *LogoutInterceptor) After(tx *sql.Tx, context map[string]any, results *any, allResults any) error {
	// params := context["params"].(map[string]any)
	if val, ok := context["session_id"].(string); ok {
		delete(Sessions, val)
	}
	return nil
}
