package plugins

import (
	"database/sql"
	"encoding/json"
	"log"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/elgs/wsl"
)

type LoginInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *LoginInterceptor) Before(tx *sql.Tx, script *string, params map[string]string, headers map[string]string) error {
	params["case"] = "lower"
	return nil
}

func (this *LoginInterceptor) After(tx *sql.Tx, result *[]interface{}) error {
	if v, ok := (*result)[3].([]map[string]string); ok {
		if len(v) == 0 {
			log.Println("Login failed.")
		} else {
			log.Println("Login succeeded.")
			loginData, err := json.Marshal(v[0])
			if err != nil {
				return nil
			}
			token, err := createJwtToken(string(loginData))
			if err != nil {
				return nil
			}
			*result = append(*result, token)
		}
	} else {
		log.Println("Login failed.")
		return nil
	}
	return nil
}

func createJwtToken(payload string) (string, error) {
	key := []byte("Some secret password")
	token, err := jose.Sign(payload, jose.HS256, key)
	return token, err
}
