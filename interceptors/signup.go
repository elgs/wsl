package interceptors

import (
	"database/sql"

	"github.com/elgs/gostrgen"
	"github.com/elgs/wsl"
)

type SignupInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *SignupInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	signupCode, err := gostrgen.RandGen(6, gostrgen.Digit, "", "")
	if err != nil {
		return err
	}
	params := context["params"].(map[string]interface{})
	params["__signup"] = signupCode

	return nil
}
