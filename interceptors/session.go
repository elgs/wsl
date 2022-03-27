package interceptors

import (
	"database/sql"

	"github.com/elgs/wsl"
)

type SessionInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *SessionInterceptor) After(tx *sql.Tx, context map[string]any, results *any, allResults any) error {
	*results = context["session"]
	return nil
}
