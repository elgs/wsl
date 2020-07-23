package interceptors

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/elgs/gosqljson"
	"github.com/elgs/wsl"
)

var sessionQuery = `
SELECT
USER.ID AS USER_ID,
USER.USERNAME,
USER.EMAIL,
USER.USER_FLAG,
USER.MODE,
USER.TIME_CREATED,
USER_SESSION.ID AS SESSION_ID,
USER_SESSION.LOGIN_TIME,
USER_SESSION.LOGIN_IP
FROM USER INNER JOIN USER_SESSION ON USER.ID=USER_SESSION.USER_ID 
WHERE USER_SESSION.ID=?`

var updateLastSeenQuery = `
UPDATE USER_SESSION 
SET LAST_SEEN_TIME=CONVERT_TZ(NOW(),'System','+0:0'),
LAST_SEEN_IP=?
WHERE ID=?
`

var sessions = make(map[string]map[string]string)

func (this *AuthInterceptor) getSession(tx *sql.Tx, sessionId string) (map[string]string, error) {
	if val, ok := sessions[sessionId]; ok {
		return val, nil
	}

	dbResult, err := gosqljson.QueryTxToMap(tx, "lower", sessionQuery, sessionId)
	if err != nil {
		return nil, err
	}
	if len(dbResult) != 1 {
		return nil, errors.New("session_not_found")
	}
	sessions[sessionId] = dbResult[0]
	return dbResult[0], nil
}

func (this *AuthInterceptor) updateLastSeen(db *sql.DB, sessionId string, ip string) {
	gosqljson.ExecDb(db, updateLastSeenQuery, ip, sessionId)
}

type AuthInterceptor struct {
	*wsl.DefaultInterceptor
}

func (this *AuthInterceptor) Before(tx *sql.Tx, context map[string]interface{}) error {

	params := context["params"].(map[string]interface{})

	if tokenString, ok := context["access_token"].(string); ok {

		session, err := this.getSession(tx, tokenString)
		if err != nil {
			return err
		}
		if app, ok := context["app"].(*wsl.WSL); ok {
			db := app.Databases["main"]
			if params, ok := context["params"].(map[string]interface{}); ok {
				clientIp := params["__client_ip"]
				go this.updateLastSeen(db, tokenString, clientIp.(string))
			}
		}

		params["__session_id"] = fmt.Sprintf("%v", session["session_id"])
		params["__user_id"] = fmt.Sprintf("%v", session["user_id"])
		params["__user_mode"] = fmt.Sprintf("%v", session["mode"])

		context["session_id"] = session["session_id"]
		context["session"] = session
		context["user_id"] = session["user_id"]
		context["user_mode"] = session["mode"]
	}
	return nil
}
