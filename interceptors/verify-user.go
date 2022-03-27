package interceptors

import (
	"database/sql"
	"errors"

	"github.com/elgs/wsl"
)

type VerifyUserInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *VerifyUserInterceptor) Before(tx *sql.Tx, context map[string]any) error {

	if _, ok := context["session"].(map[string]any); !ok {
		return errors.New("invalid_session")
	}

	return nil
}

func (this *VerifyUserInterceptor) After(tx *sql.Tx, context map[string]any, results *any, allResults any) error {
	if session, ok := context["session"].(map[string]any); ok {
		if flags, ok := session["flags"].(map[string]any); ok {
			delete(flags, "signup")
		}
	}
	return nil
}
