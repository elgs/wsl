package jobs

import (
	"github.com/elgs/wsl"
	"github.com/elgs/wsl/interceptors"
)

var cleanup = &wsl.Job{
	Cron: "@every 1h",
	MakeFunc: func(app *wsl.WSL) func() {
		app.Scripts["cleanup"] = `
		set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

		#expired_sessions
		SELECT ID FROM USER_SESSION WHERE LAST_SEEN_TIME < (@now_utc - INTERVAL 1 DAY);
		delete FROM USER_SESSION WHERE LAST_SEEN_TIME < (@now_utc - INTERVAL 1 DAY);

		delete FROM USER WHERE EXISTS (
			SELECT 1 FROM USER_FLAG 
			WHERE USER_FLAG.CREATED_TIME < (@now_utc - INTERVAL 1 HOUR)
			AND USER.ID=USER_FLAG.USER_ID
			AND USER_FLAG.CODE='signup'
		);

		delete FROM USER_FLAG WHERE CREATED_TIME < (@now_utc - INTERVAL 1 HOUR);
		`

		after := func(result map[string]any) {
			for _, sessions := range result {
				if ss, ok := sessions.([]map[string]string); ok {
					for _, session := range ss {
						id := session["id"]
						delete(interceptors.Sessions, id)
					}
				}
			}
		}

		return executeSQL(app.Databases["main"], app.Scripts["cleanup"], nil, nil, after)
	},
}
