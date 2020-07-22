package scripts

var logoutSql = `
-- set @now_utc := CONVERT_TZ(NOW(),'System','+0:0');

DELETE FROM USER_SESSION WHERE ID='__session_id';`
