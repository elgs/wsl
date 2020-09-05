package jobs

import (
	"github.com/elgs/wsl"
)

var cleanup = &wsl.Job{
	Cron: "@every 1h",
	MakeFunc: func(app *wsl.WSL) func() {
		app.Scripts["cleanup"] = `
		set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');
		DELETE FROM USER_SESSION WHERE LAST_SEEN_TIME < (@now_utc - INTERVAL 1 DAY);

		DELETE FROM USER WHERE EXISTS (
			SELECT 1 FROM USER_FLAG 
			WHERE USER_FLAG.CREATED_TIME < (@now_utc - INTERVAL 1 HOUR)
			AND USER.ID=USER_FLAG.USER_ID
			AND USER_FLAG.CODE='signup'
		);

		DELETE FROM USER_FLAG WHERE CREATED_TIME < (@now_utc - INTERVAL 1 HOUR);
		`
		return executeSQL(app.Databases["main"], app.Scripts["cleanup"], nil)
	},
}
